package pods

import (
	"container/heap"
	"container/list"
	"encoding/json"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/pkg/errors"
)

// ServicePod 模拟在线服务处理的Pod算法。一般情况下，一个应用程序进程通常会处理大量的请求，通常使用多线程进行处理，因此
// 若出现线程阻塞的情况，将会切换到另一个线程继续处理。理想情况下，在CPU看来，几乎所有的时间都在运行程序代码，而不是阻塞。
// 基于此假设，我们仅在提供一个影响服务请求处理长短的参数ComputeSlotRequired，含义是处理一个请求所需要的实际时间，不包含
// 因IO等阻塞的时间。
// 影响服务处理长短的，通常还有如下的因素，可以适当的建模：
// 1）切换线程时产生额外的开销，或者是完成服务的一些额外开销
// 2）因服务器负载高，而导致的Cache不命中，内存访问变慢等的额外IO开销
// 3）由于负载过高导致的激烈竞争
type ServicePod interface {
	core.PodAlgorithm

	// GetLoad 获取当前Pod的负载大小。用于让控制器负载均衡。具体的计算，可以根据当前Pod能够处理的容量，以及内存使用量决定
	// 返回的结果应该是在[0,1]区间。
	GetLoad() float64

	// DeliverRequest 让服务Pod处理服务。Pod可能会因为队列已满或者内存不够用而导致拒绝服务，需要控制器处理拒绝服务的逻辑。
	DeliverRequest(ctx *ServiceContext) error

	// ReturnUnhandledRequests 获取Pod还未处理的服务请求。通常在Pod需要退出之前使用。在Return完后，本Pod就不再处理这些请求。
	ReturnUnhandledRequests() []*ServiceContext
}

type ServiceContext struct {
	RequestId int

	// OnDone 在服务完成时通知控制器的回调函数
	OnDone func(requestId int)

	// SlotRequired 处理本请求所需要的CPU时间片长度
	// 本实现假设通常情况下这个值都比较小，通常的在线服务不会实际消耗太多的CPU时间，而是花更多时间在IO上。因此，请求处理
	// 不会跨Tick，必须在一个Tick内完成
	SlotRequired float64

	// MemRequired 处理本请求所需要的内存数量，是在正式处理时使用，而非未得到处理时的使用值。
	// 通常情况下，即使没有办法立即处理本服务，也需要提供部分内存来缓存本请求，可以取适当的比值反映此内存占用
	MemRequired int64
}

const SimServicePod = "SimServicePod"

var simServicePodFacory core.PodAlgorithmFactory = func(argJson string, pod *core.Pod) (core.PodAlgorithm, error) {
	arg := &SimServicePodArgs{}
	err := json.Unmarshal([]byte(argJson), arg)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal arg json")
	}
	return &simServicePod{
		baseMem: arg.BaseMem,
	}, nil
}

type SimServicePodArgs struct {
	// BaseMem 运行本服务所需要的基本内存量
	BaseMem int64 `json:"baseMem"`
}

type simServicePod struct {
	initialized      bool
	baseMem          int64
	lastUsedMem      int64
	lastAvailableMem int64
	lastLoad         float64
	lastMemUsage     float64
	queue            list.List
}

var _ ServicePod = &simServicePod{}

var (
	ErrDenial       = errors.New("Cannot Provide Service")
	ErrInitializing = errors.New("Pod is currently initializing, cannot provide service")
)

// serviceContextStoreRatio 存储服务所需要的最小内存。若实际分配内存不足所有服务请求的最小内存的和，将导致服务质量下降
const serviceContextStoreRatio = 0.2

// maxProbe 若CPU时间片开始不足服务一个请求时，往后探查的最多的请求数量。若超过这个数量，均没有找到可服务的请求，
// 则退出本轮计算。
// ServicePod以先来先服务的策略处理请求，如果这个值过大，可能导致大量的后到服务先处理。
const maxProbe = 10

func (p *simServicePod) GetLoad() float64 {
	if !p.initialized {
		return 0
	}

	// 取比值的最大值
	res := p.lastLoad
	if p.lastLoad < p.lastMemUsage {
		res = p.lastMemUsage
	}
	return res
}

func (p *simServicePod) DeliverRequest(ctx *ServiceContext) error {
	if !p.initialized {
		return ErrInitializing
	}
	if p.lastAvailableMem-p.lastUsedMem < int64(serviceContextStoreRatio*float64(ctx.MemRequired)) {
		return ErrDenial
	}
	p.queue.PushBack(ctx)
	return nil
}

func (p *simServicePod) ReturnUnhandledRequests() []*ServiceContext {
	res := make([]*ServiceContext, p.queue.Len())
	for i := 0; i < len(res); i++ {
		res[i] = p.queue.Front().Value.(*ServiceContext)
		p.queue.Remove(p.queue.Front())
	}
	return res
}

type freeCpu struct {
	cpuIdx     int
	startTime  float64
	timeRemain float64
}

type freeCpuList struct {
	list []*freeCpu
}

func (f *freeCpuList) Len() int {
	return len(f.list)
}

func (f *freeCpuList) Less(i, j int) bool {
	return f.list[i].startTime < f.list[j].startTime
}

func (f *freeCpuList) Swap(i, j int) {
	temp := f.list[i]
	f.list[i] = f.list[j]
	f.list[j] = temp
}

func (f *freeCpuList) Push(x interface{}) {
	f.list = append(f.list, x.(*freeCpu))
}

func (f *freeCpuList) Pop() interface{} {
	ret := f.list[len(f.list)-1]
	f.list = f.list[:len(f.list)-1]
	return ret
}

func (p *simServicePod) Tick(slot []float64, mem int) (Load float64, MemUsage int) {
	if !p.initialized {
		// 初始化返回
		p.lastUsedMem = p.baseMem
		p.lastAvailableMem = int64(mem)
		p.lastLoad = 0.5
		p.lastMemUsage = float64(p.baseMem) / float64(mem)
		p.initialized = true
		return 0.5, int(p.baseMem)
	}

	// 计算所需总内存
	memRequired := p.baseMem
	for cur := p.queue.Front(); cur != nil; cur = cur.Next() {
		ctx := cur.Value.(*ServiceContext)
		memRequired += ctx.MemRequired
	}

	// 内存过少时惩罚性增加处理时间
	slotMultiplier := memoryShortagePenalty(int64(mem), int64(serviceContextStoreRatio*float64(mem)))

	freeList := &freeCpuList{list: make([]*freeCpu, 0, len(slot))}
	// 记录每个CPU处理单个服务的内存最大值
	maxMemUsed := make([]int64, len(slot))
	heap.Init(freeList)
	for i := 0; i < len(slot); i++ {
		heap.Push(freeList, &freeCpu{
			cpuIdx:     i,
			startTime:  0,
			timeRemain: slot[i],
		})
	}

	probeFailedCount := 0
	for cur := p.queue.Front(); probeFailedCount < maxProbe && cur != nil; {
		ctx := cur.Value.(*ServiceContext)
		slotRequired := slotMultiplier * ctx.SlotRequired
		cannotUsed := make([]*freeCpu, 0, len(slot))
		var capableCpu *freeCpu
		// 寻找最早可用CPU
		for capableCpu = heap.Pop(freeList).(*freeCpu); capableCpu.timeRemain < slotRequired; capableCpu = heap.Pop(freeList).(*freeCpu) {
			cannotUsed = append(cannotUsed, capableCpu)
			if freeList.Len() == 0 {
				capableCpu = nil
				break
			}
		}
		// 放回不能用的CPU
		for _, cpu := range cannotUsed {
			heap.Push(freeList, cpu)
		}

		if capableCpu == nil {
			probeFailedCount++
			cur = cur.Next()
			continue
		} else {
			probeFailedCount = 0
		}

		capableCpu.timeRemain -= slotRequired
		capableCpu.startTime += slotRequired
		if ctx.OnDone != nil {
			ctx.OnDone(ctx.RequestId)
		}
		heap.Push(freeList, capableCpu)

		if ctx.MemRequired > maxMemUsed[capableCpu.cpuIdx] {
			maxMemUsed[capableCpu.cpuIdx] = ctx.MemRequired
		}

		// 删除已服务元素
		shouldRemove := cur
		cur = cur.Next()
		p.queue.Remove(shouldRemove)
	}

	// 计算CPU负载
	totalAvailableSlot := float64(0)
	totalUnusedSlot := float64(0)
	for i := 0; i < len(freeList.list); i++ {
		totalAvailableSlot += slot[freeList.list[i].cpuIdx]
		totalUnusedSlot += freeList.list[i].timeRemain
	}
	// 计算内存使用
	for i := 0; i < len(maxMemUsed); i++ {
		MemUsage += int(maxMemUsed[i])
	}
	for cur := p.queue.Front(); cur != nil; cur = cur.Next() {
		MemUsage += int(serviceContextStoreRatio * float64(cur.Value.(*ServiceContext).MemRequired))
	}
	MemUsage += int(p.baseMem)
	if MemUsage > mem {
		// 可以理解为超出的内存存储在了硬盘
		MemUsage = mem
	}

	Load = (totalAvailableSlot - totalUnusedSlot) / totalAvailableSlot
	// 更新缓存
	p.lastMemUsage = float64(MemUsage) / float64(mem)
	p.lastLoad = Load
	p.lastAvailableMem = int64(mem)
	p.lastUsedMem = int64(MemUsage)

	return
}

func (p *simServicePod) ResourceRequest() (cpu float64, mem int) {
	for cur := p.queue.Front(); cur != nil; cur = cur.Next() {
		ctx := cur.Value.(*ServiceContext)
		mem += int(ctx.MemRequired)
		cpu += ctx.SlotRequired
	}
	mem += int(p.baseMem)
	return
}

// memoryShortagePenalty 计算若内存不足时，导致的时间增加的惩罚。可以理解为，物理内存不足，导致部分内存需要存放到硬盘中，
// 从而导致Page Miss，产生额外的开销。这个开销表现为时间片的增加
func memoryShortagePenalty(memAvailable, memNeeded int64) float64 {
	if memNeeded <= memAvailable {
		return 1
	}
	return 1 + float64(memNeeded-memAvailable)/float64(memAvailable)
}
