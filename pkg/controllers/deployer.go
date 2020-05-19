// controllers contains all kinds of controller implementation, each with different features.
package controllers

import (
	"container/heap"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
)

// ControllerDeployer deploys controller at given tick. This controller will regard first Tick() call as tick 0. And it
// will deploy a controller at tick T if called DeployAt() with arguments (controller, T).
type ControllerDeployer interface {
	core.Controller

	// DeployAt deploy a controller to the cluster at given tick. If current tick is larger than tick count, this method
	// simply ignore the request.
	DeployAt(controller core.Controller, tick int, when DeployTime)
}

type controllerTimer struct {
	controller core.Controller
	when       DeployTime
	tick       int
}

type priorityQueue []*controllerTimer

func (p priorityQueue) Len() int {
	return len(p)
}

func (p priorityQueue) Less(i, j int) bool {
	return p[i].tick < p[j].tick
}

func (p priorityQueue) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}

func (p *priorityQueue) Push(x interface{}) {
	*p = append(*p, x.(*controllerTimer))
}

func (p *priorityQueue) Pop() interface{} {
	item := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return item
}

func NewControllerDeployer() ControllerDeployer {
	timers := priorityQueue(make([]*controllerTimer, 0, 10))
	heap.Init(&timers)
	return &controllerDeployer{
		tick:  0,
		queue: &timers,
	}
}

type controllerDeployer struct {
	sim   *core.SchedSim
	tick  int
	queue *priorityQueue
}

func (c controllerDeployer) Tick() {
	for c.queue.Len() > 0 && (*c.queue)[0].tick <= c.tick {
		item := heap.Pop(c.queue)
		timer := item.(*controllerTimer)
		if timer.when == BeforeUpdate {
			c.sim.RegisterBeforeUpdateController(timer.controller)
		} else if timer.when == AfterUpdate {
			c.sim.RegisterAfterUpdateController(timer.controller)
		}
	}

	c.tick++
}

type DeployTime string

var (
	BeforeUpdate DeployTime = "b"
	AfterUpdate  DeployTime = "a"
)

func (c controllerDeployer) DeployAt(controller core.Controller, tick int, when DeployTime) {
	if tick < c.tick {
		// ignore this request
		return
	}

	timer := &controllerTimer{
		controller: controller,
		tick:       tick,
		when:       when,
	}
	heap.Push(c.queue, timer)
}
