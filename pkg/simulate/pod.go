package simulate

import (
	v1 "k8s.io/api/core/v1"
	apimachineryv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	v1.Pod

	// CpuLimit CPU限制核数
	CpuLimit int

	// MemLimit Mem限制大小
	MemLimit int

	// 部署本Pod的DeploymentController
	Controller DeploymentController

	// 具体运行的算法
	Algorithm PodAlgorithm
}

// PodAlgorithmFactory 构造Pod核心算法的工厂方法，为了让算法能够访问Pod，工厂方法会在运行时得到Pod的实际指针，该指针用于
// 访问Pod的状态信息，从而计算出具体的逻辑。
type PodAlgorithmFactory func(pod *Pod) PodAlgorithm

func (p *Pod) Tick(slot []float64, mem int) (Load float64, MemUsage int) {
	return p.Algorithm.Tick(slot, mem)
}

func (p *Pod) ResourceRequest() (cpu int, mem int) {
	return p.Algorithm.ResourceRequest()
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
	ResourceRequest() (cpu int, mem int)
}

type PodBuilder struct {
	Name                      string
	Labels                    map[string]string
	Annotations               map[string]string
	NodeSelector              map[string]string
	NodeName                  string
	Affinity                  *v1.Affinity
	SchedulerName             string
	Toleration                []v1.Toleration
	PriorityClassName         string
	PreemptionPolicy          *v1.PreemptionPolicy
	TopologySpreadConstraints []v1.TopologySpreadConstraint

	CpuLimit         int
	MemLimit         int
	Controller       DeploymentController
	AlgorithmFactory PodAlgorithmFactory
}

func (builder *PodBuilder) Build() *Pod {
	if builder.AlgorithmFactory == nil || builder.Controller == nil {
		panic("Pod必须有AlgorithmFactory和DeploymentController")
	}

	p := &Pod{
		Pod: v1.Pod{
			TypeMeta: apimachineryv1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: apimachineryv1.ObjectMeta{
				Name:              builder.Name,
				CreationTimestamp: apimachineryv1.Time{},
				Labels:            builder.Labels,
				Annotations:       builder.Annotations,
				ClusterName:       "Simulator",
			},
			Spec: v1.PodSpec{
				NodeSelector:              builder.NodeSelector,
				NodeName:                  builder.NodeName,
				Affinity:                  builder.Affinity,
				SchedulerName:             builder.SchedulerName,
				Tolerations:               builder.Toleration,
				PriorityClassName:         builder.PriorityClassName,
				Priority:                  nil,
				PreemptionPolicy:          builder.PreemptionPolicy,
				TopologySpreadConstraints: builder.TopologySpreadConstraints,
			},
			Status: v1.PodStatus{
				Phase:             v1.PodPending,
				NominatedNodeName: "",
				StartTime:         nil,
				QOSClass:          "",
			},
		},
		CpuLimit:   builder.CpuLimit,
		MemLimit:   builder.MemLimit,
		Controller: builder.Controller,
	}

	p.Algorithm = builder.AlgorithmFactory(p)
	return p
}

type PodState string

var (
	RunningState   PodState = "Running"
	TerminateState PodState = "Terminate"
	// ErrorState 代表容器进入了不可恢复的错误状态。通常导致容器无法提供服务，并且无法自行恢复
	ErrorState PodState = "Error"
)

type PodEventType string

type PodEvent struct {
	Who  *Pod
	What PodEventType
}

const (
	PodPreemptEvent   = PodEventType("Preempt")
	PodTerminateEvent = PodEventType("Terminate")
)
