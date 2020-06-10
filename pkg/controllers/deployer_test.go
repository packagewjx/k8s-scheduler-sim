package controllers

import (
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"testing"
	"time"
)

type deployerTestSimulator struct {
	ch chan string
}

var _ core.SchedulerSimulator = &deployerTestSimulator{}

func (f *deployerTestSimulator) GetPod(name string) (*core.Pod, error) {
	panic("implement me")
}

func (f *deployerTestSimulator) GetKubernetesClient() kubernetes.Interface {
	panic("implement me")
}

func (f *deployerTestSimulator) GetInformerFactory() informers.SharedInformerFactory {
	panic("implement me")
}

func (f *deployerTestSimulator) Run() {
	panic("implement me")
}

func (f *deployerTestSimulator) RegisterBeforeUpdateController(controller core.Controller) {
	f.ch <- "before"
}

func (f *deployerTestSimulator) RegisterAfterUpdateController(controller core.Controller) {
	f.ch <- "after"
}

func (f *deployerTestSimulator) DeleteBeforeController(controller core.Controller) {
	panic("implement me")
}

func (f *deployerTestSimulator) DeleteAfterController(controller core.Controller) {
	panic("implement me")
}

type fakeController struct {
}

func (f *fakeController) Name() string {
	panic("implement me")
}

func (f *fakeController) Tick() {
	panic("implement me")
}

func TestDeployer(t *testing.T) {
	sim := &deployerTestSimulator{ch: make(chan string, 1)}

	deployer := NewControllerDeployer(sim)

	deployer.DeployAt(&fakeController{}, 10, BeforeUpdate)

	deployer.DeployAt(&fakeController{}, 100, AfterUpdate)

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
