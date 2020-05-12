package core

// Controller 广义控制器，管理资源或者模拟负载情况等。Controller应当通过kubernetes.Interface访问集群资源
type Controller interface {
	// Tick 更新控制器状态的函数。本函数应该是同步的
	Tick(sim *SchedSim)
}

// DeploymentController 模拟Kubernetes的控制器，根据其配置的模板构建Pod，然后通过Tick方法提交到本集群
// 用户可以实现本接口，以定制Pod的提交。如批处理任务中某些Pod先于另一些Pod提交，或者在线业务中，压力增大时提交更多的Pod的
// 逻辑。
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
