package controllers

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/informers"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util/fake"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"testing"
	"time"
)

func TestReplication(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	client := fake.NewFakeKubernetesInterface()
	sim := &replicationTestSimulator{
		client:  client,
		factory: informers.NewSharedInformerFactory(client),
	}
	stopCh := make(chan struct{})
	defer func() {
		stopCh <- struct{}{}
	}()

	replicaNum := 10
	startChan := make(chan bool, replicaNum)
	deleteChan := make(chan bool, replicaNum)
	sim.factory.Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			startChan <- true
		},
		UpdateFunc: nil,
		DeleteFunc: func(obj interface{}) {
			deleteChan <- true
		},
	})
	sim.factory.Start(stopCh)

	pods := make([]*v1.Pod, 0, 10)

	rc := NewReplicationController(sim, "test", replicaNum, func() *v1.Pod {
		pod := &v1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: string(uuid.NewUUID()),
			},
			Spec:   v1.PodSpec{},
			Status: v1.PodStatus{},
		}
		pods = append(pods, pod)
		return pod
	})

	// init tick
	rc.Tick()
	// deploy tick
	rc.Tick()

	for i := 0; i < replicaNum; i++ {
		select {
		case <-startChan:
		case <-time.After(time.Second):
			t.Error("no deploy")
		}
	}

	// 结束一个Pod，测试是否重启
	pods[0].Status.Phase = v1.PodSucceeded
	_, err := client.CoreV1().Pods(core.DefaultNamespace).UpdateStatus(context.TODO(), pods[0], metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("Update failed")
	}

	// 更新状态
	rc.Tick()

	select {
	case <-startChan:
	case <-time.After(time.Second):
		t.Error("no deploy for deleted pod")
	}

	// 调用Terminate
	rc.Terminate()
	for i := 0; i < replicaNum; i++ {
		rc.Tick()
		select {
		case <-deleteChan:
		case <-time.After(time.Second):
			t.Error("no delete")
		}
	}
}

type replicationTestSimulator struct {
	client  kubernetes.Interface
	factory k8sinformers.SharedInformerFactory
}

func (f *replicationTestSimulator) GetKubernetesClient() kubernetes.Interface {
	return f.client
}

func (f *replicationTestSimulator) GetInformerFactory() k8sinformers.SharedInformerFactory {
	return f.factory
}

func (f *replicationTestSimulator) Run() {

}

func (f *replicationTestSimulator) RegisterBeforeUpdateController(controller core.Controller) {

}

func (f *replicationTestSimulator) RegisterAfterUpdateController(controller core.Controller) {

}

func (f *replicationTestSimulator) DeleteBeforeController(controller core.Controller) {

}

func (f *replicationTestSimulator) DeleteAfterController(controller core.Controller) {

}
