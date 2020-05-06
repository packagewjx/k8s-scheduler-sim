package pkg

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/simulate"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"testing"
)

func TestScheduleOne(t *testing.T) {
	ctx := context.TODO()
	nodeBuilder := &simulate.NodeBuilder{
		Name:        "",
		Labels:      nil,
		Annotations: nil,
		Taints:      nil,
		CoreCount:   4,
		MemorySize:  1 << 32,
		Scheduler:   simulate.NewFairScheduler(),
	}
	simulator, _ := NewSchedulerSimulator(ctx, 1, nodeBuilder)
	pod := (&simulate.PodBuilder{
		Name:                      "test-1",
		Labels:                    nil,
		Annotations:               nil,
		NodeSelector:              nil,
		NodeName:                  "",
		Affinity:                  nil,
		SchedulerName:             "",
		Toleration:                nil,
		PriorityClassName:         "",
		PreemptionPolicy:          nil,
		TopologySpreadConstraints: nil,
		CpuLimit:                  1,
		MemLimit:                  1000,
		Controller:                &mockDeploymentController{},
		AlgorithmFactory:          simulate.BatchPodFactory(1000, 1000),
	}).Build()

	state := framework.NewCycleState()
	scheduleResult, err := simulator.Scheduler.Algorithm.Schedule(ctx, nil, state, &pod.Pod)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(scheduleResult)
}

type mockDeploymentController struct {
}

func (dc *mockDeploymentController) Tick() (addPod []*simulate.Pod, removePod []*simulate.Pod) {
	return []*simulate.Pod{}, []*simulate.Pod{}
}

func (dc *mockDeploymentController) InformPodEvent(event *simulate.PodEvent) {

}
