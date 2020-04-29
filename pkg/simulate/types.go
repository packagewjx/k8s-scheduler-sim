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
	Pod *Pod
	// 执行的时间片大小
	Slot float64
}
