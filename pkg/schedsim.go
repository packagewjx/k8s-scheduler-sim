package pkg

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/mock"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/simulate"
	"github.com/pkg/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/scheduler"
)

type SchedSim struct {
	client    kubernetes.Interface
	nodes     []*simulate.Node
	scheduler *scheduler.Scheduler
}

func NewSchedulerSimulator(ctx context.Context, nodeCount int) (*SchedSim, error) {
	sched, err := buildScheduler(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error building scheduler")
	}

	return &SchedSim{
		client:    mock.NewClient(),
		scheduler: sched,
	}, nil
}

func buildScheduler(ctx context.Context) (*scheduler.Scheduler, error) {
	client := mock.NewClient()
	return scheduler.New(client, informers.NewSharedInformerFactory(client, 0), scheduler.NewPodInformer(client, 0), mock.SimRecorderFactory, ctx.Done())
}

func (sim *SchedSim) Run() {

}
