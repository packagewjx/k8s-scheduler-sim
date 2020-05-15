package core

func NewReplicationController(replicaNum int) ControllerFactory {
	return func(sim *SchedSim) Controller {
		return &replicationController{
			replicaNum: replicaNum,
		}
	}
}

type ReplicationController interface {
	Controller
	SetReplicaNum(num int)
	Terminate()
}

type replicationController struct {
	sim        *SchedSim
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
