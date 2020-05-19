package core

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
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

func BuildPod(name string, cpuLimit float64, memLimit int, algorithm string, deploymentController string, initState interface{}, schedulerName string) (*v1.Pod, error) {
	stateBytes, err := json.Marshal(initState)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal initState")
	}

	podStateJson := string(stateBytes)

	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			UID:  uuid.NewUUID(),
			Annotations: map[string]string{
				PodAnnotationCpuLimit:             fmt.Sprintf("%.3f", cpuLimit),
				PodAnnotationMemLimit:             fmt.Sprintf("%d", memLimit),
				PodAnnotationAlgorithm:            algorithm,
				PodAnnotationDeploymentController: deploymentController,
				PodAnnotationInitialState:         podStateJson,
			},
		},
		// inorder to go to unscheduled queue
		Spec: v1.PodSpec{
			SchedulerName: schedulerName,
		},
		Status: v1.PodStatus{
			Phase: v1.PodPending,
		},
	}, nil
}

type PodAlgorithm interface {
	// 返回当前一个时钟滴答内的使用率。
	// slot 为Pod分配到的时间片的大小的数组，每个值取值0～1，数组长度代表分配到的CPU数量。在负载一定的情况下，
	//      若slot小，且CPU少，则可能会有更长的高压力时间。若长度为0，则代表本周期没有被调度。
	// mem Pod分配到的内存的大小，为实际的大小。Pod不能使用超过此值的内存，否则最多分配此值。
	// Load 代表单位时间内的负载压力指示，取值0～1。是使用了时间片的比值，1代表使用了所有的时间片。
	// MemUsage 实际占用的内存大小，不能超过mem
	Tick(slot []float64, mem int64) (Load float64, MemUsage int64)

	// ResourceRequest 返回本Pod在下一个周期所需要使用的CPU数和内存数。
	// cpu 节点CPU数量。注意不能超过本Pod的限制。
	// mem 节点空闲内存数量。这部分内存尚未使用，而另一部分内存则被Pod占用。Pod可以选择提高Mem，也可以降低Mem使用。
	//     不能超过本Pod的限制。
	ResourceRequest() (cpu float64, mem int64)
}

type PodAlgorithmFactory func(argJson string, pod *Pod) (PodAlgorithm, error)

var podAlgorithmMap = map[string]PodAlgorithmFactory{
	testPod: testPodAlgorithmFactory,
}

func RegisterPodAlgorithmFactory(name string, factory PodAlgorithmFactory) {
	podAlgorithmMap[name] = factory
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

// Controller 广义控制器，管理资源或者模拟负载情况等。Controller应当通过kubernetes.Interface访问集群资源
type Controller interface {
	// Tick 更新控制器状态的函数。本函数应该是同步的
	Tick()
}

// DeploymentController 模拟Kubernetes的控制器，根据其配置的模板构建Pod，然后通过Tick方法提交到本集群
// 用户可以实现本接口，以定制Pod的提交。如批处理任务中某些Pod先于另一些Pod提交，或者在线业务中，压力增大时提交更多的Pod的
// 逻辑。
// 由于部署并非100%成功，控制器需要监听集群事件，若发现部署失败，则根据情况重新部署。
type DeploymentController interface {
	Controller
}

// ServiceController 模拟Kubernetes中进行负载均衡的Service，负责接收外部请求，并分配到各个Service Pod上。
type ServiceController interface {
	Controller
}

// RunEntity 是在一个CPU队列上可执行的任务的结构，
type RunEntity struct {
	// 可被调度的Pod
	Pod *Pod
	// 执行的时间片大小
	Slot float64
}

const testPod = "test"

type testPodAlgorithm struct {
}

func (t *testPodAlgorithm) Tick(_ []float64, _ int64) (Load float64, MemUsage int64) {
	return 1, 1
}

func (t *testPodAlgorithm) ResourceRequest() (cpu float64, mem int64) {
	return 1, 1
}

var testPodAlgorithmFactory PodAlgorithmFactory = func(argJson string, pod *Pod) (PodAlgorithm, error) {
	return &testPodAlgorithm{}, nil
}
