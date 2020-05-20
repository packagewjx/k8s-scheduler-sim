package controllers

import (
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"testing"
	"time"
)

type fakeSchedulerSimulator struct {
	ch chan string
}

func (f *fakeSchedulerSimulator) GetKubernetesClient() kubernetes.Interface {
	panic("implement me")
}

func (f *fakeSchedulerSimulator) GetInformerFactory() informers.SharedInformerFactory {
	panic("implement me")
}

func (f *fakeSchedulerSimulator) Run() {
	panic("implement me")
}

func (f *fakeSchedulerSimulator) RegisterBeforeUpdateController(controller core.Controller) {
	f.ch <- "before"
}

func (f *fakeSchedulerSimulator) RegisterAfterUpdateController(controller core.Controller) {
	f.ch <- "after"
}

func (f *fakeSchedulerSimulator) DeleteBeforeController(controller core.Controller) {
	panic("implement me")
}

func (f *fakeSchedulerSimulator) DeleteAfterController(controller core.Controller) {
	panic("implement me")
}

type fakeController struct {
}

func (f fakeController) Tick() {
	panic("implement me")
}

func TestDeployer(t *testing.T) {
	sim := &fakeSchedulerSimulator{ch: make(chan string, 1)}

	deployer := NewControllerDeployer(sim)

	deployer.DeployAt(fakeController{}, 10, BeforeUpdate)

	deployer.DeployAt(fakeController{}, 100, AfterUpdate)

	for i := 0; i < 10; i++ {
		deployer.Tick()
	}

	select {
	case s := <-sim.ch:
		if s != "before" {
			t.Error("error")
		}
	case <-time.After(time.Second):
		t.Error("error")
	}

	for i := 0; i < 90; i++ {
		deployer.Tick()
	}
	select {
	case s := <-sim.ch:
		if s != "after" {
			t.Error("error")
		}
	case <-time.After(time.Second):
		t.Error("error")
	}
}
