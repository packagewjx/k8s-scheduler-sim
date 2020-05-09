package simulate

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Pod struct {
	v1.Pod

	// CpuLimit CPU限制核数。小数，用于控制最多使用多少个核。小数部分代表能够一个CPU时间片的比例。
	CpuLimit float64

	// MemLimit Mem限制大小，单位为字节
	MemLimit int64

	// 具体运行的算法
	Algorithm PodAlgorithm
}

const (
	PodAnnotationCpuLimit             = "github.com/packagewjx/cpulimit"
	PodAnnotationMemLimit             = "github.com/packagewjx/memlimit"
	PodAnnotationAlgorithm            = "github.com/packagewjx/algorithm"
	PodAnnotationInitialState         = "github.com/packagewjx/initstate"
	PodAnnotationDeploymentController = "github.com/packagewjx/deploymentcontroller"
)

func (p *Pod) DeepCopyObject() runtime.Object {
	corePodClone := p.Pod.DeepCopy()
	return &Pod{
		Pod:       *corePodClone,
		CpuLimit:  p.CpuLimit,
		MemLimit:  p.MemLimit,
		Algorithm: p.Algorithm,
	}
}

type PodAlgorithm interface {
	// 返回当前一个时钟滴答内的使用率。
	// slot 为Pod分配到的时间片的大小的数组，每个值取值0～1，数组长度代表分配到的CPU数量。在负载一定的情况下，
	//      若slot小，且CPU少，则可能会有更长的高压力时间。若长度为0，则代表本周期没有被调度。
	// mem Pod分配到的内存的大小，为实际的大小。Pod不能使用超过此值的内存，否则最多分配此值。
	// Load 代表单位时间内的负载压力指示，取值0～1。是使用了时间片的比值，1代表使用了所有的时间片。
	// MemUsage 实际占用的内存大小，不能超过mem
	Tick(slot []float64, mem int) (Load float64, MemUsage int)

	// ResourceRequest 返回本Pod在下一个周期所需要使用的CPU数和内存数。
	// cpu 节点CPU数量。注意不能超过本Pod的限制。
	// mem 节点空闲内存数量。这部分内存尚未使用，而另一部分内存则被Pod占用。Pod可以选择提高Mem，也可以降低Mem使用。
	//     不能超过本Pod的限制。
	ResourceRequest() (cpu float64, mem int)
}

type PodAlgorithmFactory func(stateJson string, pod *Pod) (PodAlgorithm, error)

var podAlgorithmMap = map[string]PodAlgorithmFactory{
	BatchPodName: BatchPodFactory,
}

func GetPodAlgorithmFactory(name string) (factory PodAlgorithmFactory, exist bool) {
	factory, exist = podAlgorithmMap[name]
	return
}

type PodEventType string

type PodEvent struct {
	Who  *Pod
	What PodEventType
}

const (
	PodPreemptEvent   = PodEventType("Preempt")
	PodTerminateEvent = PodEventType("Terminate")
)
