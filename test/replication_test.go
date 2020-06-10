package test

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/controllers"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/pods"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"testing"
)

func TestReplicationController(t *testing.T) {
	sim := core.NewSchedulerSimulator(1000)

	rc := controllers.NewReplicationController(sim, "rep", 10, func() *v1.Pod {
		uid := uuid.NewUUID()
		return &v1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					core.PodAnnotationCpuLimit:     "1",
					core.PodAnnotationMemLimit:     "1",
					core.PodAnnotationInitialState: "",
					core.PodAnnotationAlgorithm:    pods.SimServicePod,
				},
				UID:  uid,
				Name: string(uid),
			},
			Spec: v1.PodSpec{
				SchedulerName: v1.DefaultSchedulerName,
			},
			Status: v1.PodStatus{},
		}
	})
	sim.RegisterBeforeUpdateController(rc)

	node := core.BuildNode("test", "10", "10000", "1000", core.FairScheduler)
	_, err := sim.GetKubernetesClient().CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error create node: %v", err)
	}

	sim.Run()
}
