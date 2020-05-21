package controllers

import (
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
	logrus.SetLevel(logrus.DebugLevel)
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
	sim.factory.Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			startChan <- true
		},
		UpdateFunc: nil,
		DeleteFunc: nil,
	})
	sim.factory.Start(stopCh)

	rc := NewReplicationController(sim, "test", replicaNum, func() *v1.Pod {
		return &v1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: string(uuid.NewUUID()),
			},
			Spec:   v1.PodSpec{},
			Status: v1.PodStatus{},
		}
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
