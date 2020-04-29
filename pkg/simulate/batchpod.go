package simulate

// BatchPod 模拟批处理任务的Pod，默认一直在跑任务，因此负载一直为1，内存使用基本固定
type BatchPod struct {
	*BasePod
	memUsage  float64
	totalTick float64
}

func (pod *BatchPod) ResourceRequest() (cpu int, mem float64) {
	return pod.cpuLimit, pod.memLimit
}

func (pod *BatchPod) GetState() PodState {
	if pod.totalTick < 0 {
		return TerminateState
	} else {
		return RunningState
	}
}

func NewBatchPod(name string, priority int, controller DeploymentController, memUsage float64, totalTick float64) Pod {
	return &BatchPod{
		BasePod: &BasePod{
			name:                 name,
			priority:             priority,
			deploymentController: controller,
		},
		memUsage:  memUsage,
		totalTick: totalTick,
	}
}

func (pod *BatchPod) Tick(slot []float64, mem float64) (Load, MemUsage float64) {
	if pod.totalTick < 0 {
		return 0, 0
	}
	slotSum := float64(0)
	for i := 0; i < len(slot); i++ {
		slotSum += slot[i]
	}

	// 内存给多少用多少，但是不超过定义的memUsage
	memUsage := pod.memUsage
	if memUsage > mem {
		memUsage = mem

		// 由于内存不太够，理论上会影响执行速度，因此修改slotSum
		slotSum *= mem / pod.memUsage
	}
	pod.totalTick -= slotSum
	return 1, memUsage
}
