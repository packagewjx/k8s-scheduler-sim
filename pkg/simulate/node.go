package simulate

import (
	"github.com/packagewjx/k8s-scheduler-sim/pkg/metrics"
)

type Node struct {
	// CPU核数
	CoreCount int
	// 内存大小，以字节为单位，应该是标准化了的
	MemorySize float64
	// 调度器，用于分配CPU时间片
	Scheduler CoreScheduler
	// 本节点上的Pod
	Pods map[string]Pod
	// CpuState 反映当前CPU状态，每个CPU上有自己的RunEntity队列，表示在上一个周期中运行的所有Pod
	CpuState [][]*RunEntity

	// 上一轮的CPU使用百分比，用于查看是否有资源竞争，模拟高资源竞争时CPU处理能力的下降
	lastCpuUsage float64
}

func NewNode(coreCount int, memorySize float64, scheduler CoreScheduler) *Node {
	return &Node{
		CoreCount:    coreCount,
		CpuState:     make([][]*RunEntity, coreCount),
		MemorySize:   memorySize,
		Scheduler:    scheduler,
		Pods:         make(map[string]Pod),
		lastCpuUsage: 0,
	}
}

// AddPod 将Pod调度到本节点运行。新的Pod将会在下一次Tick的时候分配时间片运行
func (n *Node) AddPod(pod Pod) {
	n.Pods[pod.Name()] = pod
}

// 根据节点拥有的Pod，更新当前的节点状态，包括资源使用率，Pod状态等
// TODO 内存分配
func (n *Node) Tick() *metrics.TickMetrics {
	terminatedPods := make([]Pod, 0)
	type PodResource struct {
		slot     []float64
		mem      float64
		load     float64
		memUsage float64
	}

	readyPods := make([]Pod, 0, len(n.Pods))
	podResource := make([]*PodResource, 0, len(n.Pods))
	podIdxMap := make(map[Pod]int)

	// 查看Pod的状态
	for _, pod := range n.Pods {
		if pod.GetState() == RunningState {
			_, mem := pod.ResourceRequest()
			podIdxMap[pod] = len(podResource)
			readyPods = append(readyPods, pod)
			podResource = append(podResource, &PodResource{
				slot: make([]float64, 0),
				mem:  mem,
			})
		} else {
			terminatedPods = append(terminatedPods, pod)
		}
	}

	for i := 0; i < len(terminatedPods); i++ {
		// 通知控制器停止了的Pod
		terminatedPods[i].DeploymentController().InformPodEvent(&PodEvent{
			Who:  terminatedPods[i],
			What: PodTerminateEvent,
		})
		// 从本节点移除
		delete(n.Pods, terminatedPods[i].Name())
	}

	// 执行调度算法
	cpuState := n.Scheduler.Schedule(readyPods, n.CpuState)
	for i := 0; i < len(cpuState); i++ {
		for j := 0; j < len(cpuState[i]); j++ {
			entity := cpuState[i][j]
			podIdx := podIdxMap[entity.Pod]
			podResource[podIdx].slot = append(podResource[podIdx].slot, entity.Slot)
		}
	}

	// 根据分配结果更新Pod的执行状态
	memUsed := float64(0)
	load := float64(0)
	for i := 0; i < len(readyPods); i++ {
		cpuPressureReduction(podResource[i].slot, n.lastCpuUsage)
		podResource[i].load, podResource[i].memUsage = readyPods[i].Tick(podResource[i].slot, podResource[i].mem)

		// 计算统计
		memUsed += podResource[i].memUsage
		load += podResource[i].load
	}

	// 根据Pod返回的负载信息，计算CPU统计数据
	cpuUsed := float64(0)
	for i := 0; i < n.CoreCount; i++ {
		for j := 0; j < len(cpuState[i]); j++ {
			entity := cpuState[i][j]
			podIdx := podIdxMap[entity.Pod]
			cpuUsed += entity.Slot * podResource[podIdx].load
		}
	}

	cpuUsage := cpuUsed / float64(n.CoreCount)
	n.lastCpuUsage = cpuUsage

	return &metrics.TickMetrics{
		CpuUsage: cpuUsage,
		MemUsage: memUsed / n.MemorySize,
		Load:     load,
	}

}

// 当高CPU压力时，减少实际slot数，反映CPU速率的下降。注意在最后统计时，使用没有减少的slot。
func cpuPressureReduction(slot []float64, cpuUsage float64) {
	// 只在大于0.7时触发
	if cpuUsage >= 0.7 {
		// 线性变化，减少时间片
		reduction := 1 - (cpuUsage - 0.7)

		for i := 0; i < len(slot); i++ {
			slot[i] *= reduction
		}
	}
}
