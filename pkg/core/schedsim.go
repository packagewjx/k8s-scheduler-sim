package core

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/informers"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/metrics"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/mock"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/scheduler"
	"time"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const DefaultNamespace = ""

type SchedSim struct {
	Client                kubernetes.Interface
	Nodes                 cache.Store
	DeploymentControllers cache.Store
	PriorityClasses       cache.Store
	Pods                  cache.Store
	Scheduler             *scheduler.Scheduler
	InformerFactory       k8sinformers.SharedInformerFactory
	cancelFunc            context.CancelFunc

	// 当前的时钟周期数
	Tick int
	// 总运行时钟周期数
	TotalTick int

	// BeforeUpdate 在更新Pod状态与Node状态之前调用的控制器函数
	BeforeUpdate []Controller

	// AfterUpdate 在更新Pod状态之后调用的控制器函数，通常用于监控统计等
	AfterUpdate []Controller

	// 用于控制调度器调度添加速率的两个通道
	addCount  int
	bindPodCh chan string
}

var (
	PodKeyFunc cache.KeyFunc = func(obj interface{}) (string, error) {
		if pod, ok := obj.(*Pod); ok {
			return pod.Name, nil
		} else if pod, ok := obj.(*v1.Pod); ok {
			return pod.Name, nil
		} else if pod, ok := obj.(*v1.Pod); ok {
			return pod.Name, nil
		} else {
			return "", fmt.Errorf("error getting key from %v", obj)
		}
	}
	NodeKeyFunc cache.KeyFunc = func(obj interface{}) (string, error) {
		if node, ok := obj.(*Node); ok {
			return node.Name, nil
		} else if node, ok := obj.(*v1.Node); ok {
			return node.Name, nil
		} else {
			return "", fmt.Errorf("error getting key from %v", obj)
		}
	}
	PriorityClassKeyFucn cache.KeyFunc = func(obj interface{}) (string, error) {
		if cls, ok := obj.(*schedulingv1.PriorityClass); ok {
			return cls.Name, nil
		} else {
			return "", fmt.Errorf("invalid type, not PriorityClass")
		}
	}
)

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		ForceFormatting: true,
	})
}

// NewSchedulerSimulator 创建一个新的集群。totalTick为模拟集群的总运行周期。
func NewSchedulerSimulator(totalTick int) *SchedSim {
	rootCtx, cancel := context.WithCancel(context.Background())
	sim := &SchedSim{
		Client:                nil,
		Nodes:                 cache.NewStore(NodeKeyFunc),
		DeploymentControllers: nil,
		PriorityClasses:       nil,
		Pods:                  cache.NewStore(PodKeyFunc),
		Scheduler:             nil,
		TotalTick:             totalTick,
		cancelFunc:            cancel,
		bindPodCh:             make(chan string),
	}

	client, err := NewClient(sim)
	if err != nil {
		panic(fmt.Sprintf("error create client: %s", err))
	}
	sim.Client = client
	sim.InformerFactory = informers.NewSharedInformerFactory(client)
	// explicitly trigger the creation of these informers, and then start the factory to let the informer subscribe
	sim.InformerFactory.Core().V1().Nodes().Informer()
	sim.InformerFactory.Core().V1().Pods().Informer()
	sim.InformerFactory.Start(rootCtx.Done())
	<-time.After(10 * time.Millisecond) // ensure informer topic subscription.

	sched, err := buildScheduler(rootCtx, sim.InformerFactory, client)
	if err != nil {
		panic(err)
	}
	sim.Scheduler = sched
	go sim.Scheduler.Run(rootCtx)

	return sim
}

func buildScheduler(ctx context.Context, factory k8sinformers.SharedInformerFactory, client kubernetes.Interface) (*scheduler.Scheduler, error) {
	podInformer := factory.Core().V1().Pods()
	return scheduler.New(client, factory, podInformer, mock.SimRecorderFactory, ctx.Done())
}

func (sim *SchedSim) Run() {
	defer sim.cancelFunc()

	nodeMetrics := make(map[*Node]metrics.Aggregator)

	for tick := 0; tick < sim.TotalTick; tick++ {
		logrus.Infof("Tick %d", tick)
		logrus.Debug("Running BeforeUpdate Controllers")

		for _, controller := range sim.BeforeUpdate {
			controller.Tick(sim)
		}
		//等待bind
		shouldBreak := false
		timeoutCh := time.After(time.Second)
		completed := 0
		for i := 0; i < sim.addCount && !shouldBreak; i++ {
			logrus.Debugf("Waiting to schedule pod, %d remaining", sim.addCount)
			select {
			case res := <-sim.bindPodCh:
				completed++
				podName := res[1:]
				if res[0] == 'T' {
					logrus.Debugf("Pod %s bind success", podName)
				} else /*F*/ {
					logrus.Infof("Pod %s scheduled failed", podName)
				}
			case <-timeoutCh:
				// 若调度超时则退出
				logrus.Infof("Schedule time out, %d remaining pod", sim.addCount)
				shouldBreak = true
			}
		}
		sim.addCount -= completed

		logrus.Debug("Updating Node status")
		nodes := sim.Nodes.List()
		currentMetrics := make([]*metrics.PeriodMetrics, 0, len(nodes))
		for _, item := range nodes {
			node := item.(*Node)
			logrus.Debugf("Updating Node %s", node.Name)
			met := node.Tick(sim.Client)
			aggregator, ok := nodeMetrics[node]
			if !ok {
				aggregator = metrics.NewAggregator()
				nodeMetrics[node] = aggregator
			}
			periodMetrics := aggregator.Aggregate(met)
			currentMetrics = append(currentMetrics, periodMetrics)
		}

		// 显示各个节点的状态
		// 打印表头
		fmt.Println("Node\tCPU  \tCPUALL\tCPU60\tCPU300\tCPU1500\tMem  \tMemALL\tMem60\tMem300\tMem1500\tLoad \tLoadALL\tLoad60\tLoad300\tLoad1500")
		for i, metric := range currentMetrics {
			fmt.Printf("%s\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\t%.3f\n", nodes[i].(*Node).Name,
				metric.CpuUsageLastTick, metric.CpuUsageAverage, metric.CpuUsageAverageIn60Ticks, metric.CpuUsageAverageIn300Ticks, metric.CpuUsageAverageIn1500Ticks,
				metric.MemUsageLastTick, metric.MemUsageAverage, metric.MemUsageAverageIn60Ticks, metric.MemUsageAverageIn300Ticks, metric.MemUsageAverageIn1500Ticks,
				metric.LoadLastTick, metric.LoadAverage, metric.LoadAverageIn60Ticks, metric.LoadAverageIn300Ticks, metric.LoadAverageIn1500Ticks)
		}

		logrus.Debug("Running AfterUpdate Controllers")
		// 运行后更新控制器
		for _, controller := range sim.AfterUpdate {
			controller.Tick(sim)
		}
	}

}

type controllerTiming int

var (
	beforeUpdate = controllerTiming(1)
	afterUpdate  = controllerTiming(2)
)

func (sim *SchedSim) RegisterBeforeUpdateController(controller Controller) {
	sim.registerController(controller, beforeUpdate)
}

func (sim *SchedSim) RegisterAfterUpdateController(controller Controller) {
	sim.registerController(controller, afterUpdate)
}

func (sim *SchedSim) DeleteBeforeController(controller Controller) {
	sim.deleteController(controller, beforeUpdate)
}

func (sim *SchedSim) DeleteAfterController(controller Controller) {
	sim.deleteController(controller, afterUpdate)
}

func (sim *SchedSim) registerController(controller Controller, timing controllerTiming) {
	switch timing {
	case beforeUpdate:
		sim.BeforeUpdate = append(sim.BeforeUpdate, controller)
	case afterUpdate:
		sim.AfterUpdate = append(sim.AfterUpdate, controller)
	}
}

func (sim *SchedSim) deleteController(controller Controller, timing controllerTiming) {
	var arr *[]Controller
	switch timing {
	case beforeUpdate:
		arr = &sim.BeforeUpdate
	case afterUpdate:
		arr = &sim.AfterUpdate
	default:
		panic("Invalid argument")
	}

	idx := 0
	for ; idx < len(*arr) && (*arr)[idx] != controller; idx++ {
	}

	if idx < len(*arr) {
		(*arr)[idx] = (*arr)[len(*arr)-1]
		*arr = (*arr)[:len(*arr)-1]
	}
}

// podAdded 通知创建了新的Pod。注意，本函数应该运行在与SchedSim不同的Goroutine中，否则会永久阻塞。‘
// 由于通道无法保证完全的同步，因此使用本方法同步的通知
func (sim *SchedSim) podAdded(podName string) {
	sim.addCount++
}

// podScheduledFailed 通知调度成功。注意，本函数应该运行在与SchedSim不同的Goroutine中，否则会永久阻塞。
func (sim *SchedSim) podScheduledSuccess(podName string) {
	sim.bindPodCh <- "T" + podName
}

// podScheduledFailed 通知调度失败。注意，本函数应该运行在与SchedSim不同的Goroutine中，否则会永久阻塞。
func (sim *SchedSim) podScheduledFailed(podName string) {
	sim.bindPodCh <- "F" + podName
}
