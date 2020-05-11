package pkg

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/informers"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/mock"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/simulate"
	v1 "k8s.io/api/core/v1"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/scheduler"
	"time"
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
}

var (
	podKeyFunc cache.KeyFunc = func(obj interface{}) (string, error) {
		if pod, ok := obj.(*simulate.Pod); ok {
			return pod.Name, nil
		} else if pod, ok := obj.(*v1.Pod); ok {
			return pod.Name, nil
		} else if pod, ok := obj.(*v1.Pod); ok {
			return pod.Name, nil
		} else {
			return "", fmt.Errorf("error getting key from %v", obj)
		}
	}
	nodeKeyFunc cache.KeyFunc = func(obj interface{}) (string, error) {
		if node, ok := obj.(*simulate.Node); ok {
			return node.Name, nil
		} else if node, ok := obj.(*v1.Node); ok {
			return node.Name, nil
		} else {
			return "", fmt.Errorf("error getting key from %v", obj)
		}
	}
)

func NewSchedulerSimulator() *SchedSim {
	rootCtx, cancel := context.WithCancel(context.Background())
	sim := &SchedSim{
		Client:                nil,
		Nodes:                 cache.NewStore(nodeKeyFunc),
		DeploymentControllers: nil,
		PriorityClasses:       nil,
		Pods:                  cache.NewStore(podKeyFunc),
		Scheduler:             nil,
		cancelFunc:            cancel,
	}

	client := &simClient{sim: sim}
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

	return sim
}

func buildScheduler(ctx context.Context, factory k8sinformers.SharedInformerFactory, client kubernetes.Interface) (*scheduler.Scheduler, error) {
	podInformer := factory.Core().V1().Pods()
	return scheduler.New(client, factory, podInformer, mock.SimRecorderFactory, ctx.Done())
}

func (sim *SchedSim) Run() {
	defer sim.cancelFunc()
}
