package simulate

// DeploymentController 模拟Kubernetes的控制器，根据其配置的模板构建Pod，然后通过Tick方法提交到本集群
// 用户可以实现本接口，以定制Pod的提交。如批处理任务中某些Pod先于另一些Pod提交，或者在线业务中，压力增大时提交更多的Pod的
// 逻辑。
type DeploymentController interface {
	// Tick 返回当前时钟周期所要部署的Pod和删除的Pod
	Tick() (addPod []Pod, removePod []Pod)

	// InformPodEvent 当Pod状态改变后，通知本Controller。决定是否在下一轮Tick时采取行动，如重新部署被抢占或错误的Pod
	InformPodEvent(event *PodEvent)
}

// ServiceController 模拟Kubernetes中进行负载均衡的Service，负责接收外部请求，并分配到各个Service Pod上。
type ServiceController interface {
}

// RunEntity 是在一个CPU队列上可执行的任务的结构，
type RunEntity struct {
	// 可被调度的Pod
	Pod Pod
	// 执行的时间片大小
	Slot float64
}

// Pod 模拟一个Pod
type Pod interface {
	// Name 返回本Pod的名称，应该唯一，否则会出错
	Name() string

	// Priority Pod的优先级。数值越高，则抢占CPU时更有机会占用
	Priority() int

	// Type Pod的类型，调度器可以根据此使用不同的调度策略
	Type() string

	// ResourceLimit 返回本Pod的CPU限制使用量，内存的限制使用量
	ResourceLimit() (cpuLimit int, memLimit float64)

	// ResourceRequest 返回本Pod在下一个周期所需要使用的CPU数和内存数。
	// nodeCpu 节点CPU数量。注意不能超过本Pod的限制。
	// nodeMem 节点空闲内存数量。这部分内存尚未使用，而另一部分内存则被Pod占用。Pod可以选择提高Mem，也可以降低Mem使用。
	// 不能超过本Pod的限制。
	ResourceRequest() (cpu int, mem float64)

	// 返回当前一个时钟滴答内的使用率。
	// slot 为Pod分配到的时间片的大小的数组，每个值取值0～1，数组长度代表分配到的CPU数量。在负载一定的情况下，
	//       若slot小，且CPU少，则可能会有更长的高压力时间。若长度为0，则代表本周期没有被调度。
	// mem Pod分配到的内存的大小，为实际的大小。Pod不能使用超过此值的内存，否则最多分配此值。
	// Load 代表单位时间内的负载压力指示，取值0～1。是使用了时间片的比值，1代表使用了所有的时间片。
	// MemUsage 实际占用的内存大小，不能超过mem
	Tick(slot []float64, mem float64) (Load, MemUsage float64)

	// GetState 返回当前Pod的状态
	GetState() PodState

	// DeploymentController 返回部署本Pod的Deployment控制器
	DeploymentController() DeploymentController
}

type BasePod struct {
	name                 string
	podType              string
	priority             int
	deploymentController DeploymentController
	cpuLimit             int
	memLimit             float64
}

func (pod *BasePod) GetState() PodState {
	panic("implement me")
}

func (pod *BasePod) Type() string {
	return pod.podType
}

func (pod *BasePod) ResourceRequest() (cpu int, mem float64) {
	panic("implement me")
}

func (pod *BasePod) Tick(_ []float64, _ float64) (Load, MemUsage float64) {
	panic("implement me")
}

func (pod *BasePod) Name() string {
	return pod.name
}

func (pod *BasePod) DeploymentController() DeploymentController {
	return pod.deploymentController
}

func (pod *BasePod) Priority() int {
	return pod.priority
}

func (pod *BasePod) ResourceLimit() (cpuLimit int, memLimit float64) {
	return pod.cpuLimit, pod.memLimit
}

type PodState interface {
	What() string
}

type podStateImpl struct {
	what string
}

func (p podStateImpl) What() string {
	return p.what
}

var (
	RunningState   = &podStateImpl{what: "Running"}
	TerminateState = &podStateImpl{what: "Terminate"}
	// ErrorState 代表容器进入了不可恢复的错误状态。通常导致容器无法提供服务，并且无法
	ErrorState = &podStateImpl{what: "Error"}
)

type PodEvent struct {
	Who  Pod
	What string
}

const (
	PodPreemptEvent   = "Preempt"
	PodTerminateEvent = "Terminate"
)
