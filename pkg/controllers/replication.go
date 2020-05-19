package controllers

import "github.com/packagewjx/k8s-scheduler-sim/pkg/core"

func NewReplicationController(sim core.SchedulerSimulator, replicaNum int) ReplicationController {
	return &replicationController{
		sim:        sim,
		replicaNum: replicaNum,
	}
}

type ReplicationController interface {
	core.Controller
	SetReplicaNum(num int)
	Terminate()
}

type replicationController struct {
	sim        core.SchedulerSimulator
	replicaNum int
}

func (r *replicationController) Tick() {
	panic("implement me")
}

func (r *replicationController) SetReplicaNum(num int) {
	r.replicaNum = num
}

func (r *replicationController) Terminate() {
	panic("implement me")
}
