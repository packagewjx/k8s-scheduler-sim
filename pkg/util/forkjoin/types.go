package forkjoin

// Pool 分治线程池接口。使用多线程的方式执行分而治之的任务。
type Pool interface {
	// Execute 提交一个新的ForkJoin分治任务，同步执行完成
	Execute(task Task) interface{}

	// Submit 提交一个新的ForkJoin分治任务，异步回调通知完成
	Submit(task Task, resultChan chan interface{})

	Shutdown()
}

// Task 分治算法执行接口，运行包含将数据分区以及整合的逻辑。通过调用pool.RunForkTask()的方法，将子任务提交给线程池运行。
type Task interface {
	Fork() (tasks []Task)

	IsLeaf() bool

	Leaf() interface{}

	Join(results []interface{}) interface{}
}
