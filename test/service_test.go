package test

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/controllers"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/pods"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestService(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	fmt.Printf("Main GoRoutine: %d\n", util.GetGoRoutineId())

	simulator := core.NewSchedulerSimulator(1000)

	node := core.BuildNode("test-node", "20", "10G", "1000", core.FairScheduler)
	_, err := simulator.GetKubernetesClient().CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error creating test node: %v", err)
	}

	requestId := 0
	service := controllers.NewServiceController(simulator, "test", 10, func() []*pods.ServiceContext {
		res := make([]*pods.ServiceContext, 10000)
		for i := 0; i < len(res); i++ {
			requestId++
			res[i] = &pods.ServiceContext{
				RequestId:    requestId,
				OnDone:       nil,
				SlotRequired: 0.05,
				MemRequired:  1024,
			}
		}
		return res
	}, &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				core.PodAnnotationMemLimit: "100",
				core.PodAnnotationCpuLimit: "1",
			},
		},
		Spec: v1.PodSpec{
			SchedulerName: v1.DefaultSchedulerName,
		},
		Status: v1.PodStatus{},
	})

	simulator.RegisterBeforeUpdateController(service)

	simulator.Run()
}
