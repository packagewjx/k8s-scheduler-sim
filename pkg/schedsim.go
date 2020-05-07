package pkg

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/mock"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/simulate"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
	"github.com/pkg/errors"
	v1 "k8s.io/api/scheduling/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/scheduler"
)

var DefaultNamespace = ""

type SchedSim struct {
	Client                kubernetes.Interface
	Nodes                 []*simulate.Node
	PendingPods           []*simulate.Pod
	DeploymentControllers []*simulate.DeploymentController
	PriorityClasses       []*v1.PriorityClass
	Scheduler             *scheduler.Scheduler
}

func NewSchedulerSimulator(ctx context.Context, nodeCount int, nodeBuilder *simulate.NodeBuilder) (*SchedSim, error) {
	sim := &SchedSim{}
	client := &SimClient{Sim: sim}
	sched, err := buildScheduler(ctx, client)
	if err != nil {
		return nil, errors.Wrap(err, "error building scheduler")
	}

	sim.Scheduler = sched
	sim.Client = client

	nodes := make([]*simulate.Node, nodeCount)
	err = util.GetMessageQueue().NewTopic(TopicNode)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < nodeCount; i++ {
		nodeBuilder.Name = fmt.Sprintf("SimNode-%d", i)
		nodes[i] = nodeBuilder.Build()

		event := &watch.Event{
			Type:   watch.Added,
			Object: nodes[i],
		}
		err = util.GetMessageQueue().Publish(TopicNode, event)
		if err != nil {
			fmt.Println(err)
		}
	}

	return sim, nil
}

func buildScheduler(ctx context.Context, client kubernetes.Interface) (*scheduler.Scheduler, error) {
	return scheduler.New(client, informers.NewSharedInformerFactory(client, 0), scheduler.NewPodInformer(client, 0), mock.SimRecorderFactory, ctx.Done())
}

func (sim *SchedSim) Run() {

}
