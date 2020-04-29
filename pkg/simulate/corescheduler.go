package simulate

// CoreScheduler 模拟操作系统内核的调度器，调度同一个节点上的所有可运行的Pod
type CoreScheduler interface {
	// Schedule 将readyPods调度到各个CPU中，并分配相应的时间片执行。
	// cpuState 长度与节点实际CPU数量相等，是上个周期内节点CPU所运行的所有的Pod的队列。
	//          初始化时，每个队列为空，但是有CPU核数个队列
	Schedule(readyPods []Pod, cpuState [][]*RunEntity) [][]*RunEntity
}

func NewFairScheduler() CoreScheduler {
	return &FairScheduler{}
}

// FairScheduler 是完全公平的调度器，不会理会Priority的限制，完全公平的将一个CPU分配给该所有将在该CPU上运行的Pod使用
type FairScheduler struct {
}

func (s *FairScheduler) Schedule(readyPods []Pod, cpuState [][]*RunEntity) [][]*RunEntity {
	totalCpu := len(cpuState)

	// 使用RoundRobin策略
	newState := make([][]*RunEntity, totalCpu)

	// 没有Pod运行时，返回空队列
	if len(readyPods) == 0 {
		return newState
	}
	cpuIdx := 0
	for i := 0; i < len(readyPods); i++ {
		cpu, _ := readyPods[i].ResourceRequest()
		cpuLimit, _ := readyPods[i].ResourceLimit()
		// 检查不能超过界限，也不能超过节点的CPU数
		if cpu > totalCpu {
			cpu = totalCpu
		}
		if cpu > cpuLimit {
			cpu = cpuLimit
		}
		for j := 0; j < cpu; j++ {
			newState[cpuIdx] = append(newState[cpuIdx], &RunEntity{
				Pod:  readyPods[i],
				Slot: 0,
			})

			cpuIdx++
			if cpuIdx >= totalCpu {
				cpuIdx = 0
			}
		}
	}

	// 平均时间片
	for i := 0; i < totalCpu; i++ {
		podCount := len(newState[i])
		if podCount == 0 {
			continue
		}
		slot := 1.0 / float64(podCount)
		for j := 0; j < podCount; j++ {
			newState[i][j].Slot = slot
		}
	}

	return newState
}
