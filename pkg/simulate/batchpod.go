package simulate

import v1 "k8s.io/api/core/v1"

// BatchPodAlgorithm 模拟批处理任务的Pod，默认一直在跑任务，因此负载一直为1，内存使用基本固定
type BatchPodAlgorithm struct {
	Pod       *Pod
	MemUsage  int
	TotalTick float64
}

func FactoryMethod(memUsage int, totalTick float64) PodAlgorithmFactory {
	return func(pod *Pod) PodAlgorithm {
		return &BatchPodAlgorithm{
			Pod:       pod,
			MemUsage:  memUsage,
			TotalTick: totalTick,
		}
	}
}

func (alg *BatchPodAlgorithm) ResourceRequest() (cpu int, mem int) {
	return alg.Pod.CpuLimit, alg.MemUsage
}

func (alg *BatchPodAlgorithm) Tick(slot []float64, mem int) (Load float64, MemUsage int) {
	if alg.TotalTick < 0 {
		alg.Pod.Status.Phase = v1.PodSucceeded
		return 0, 0
	}

	slotSum := float64(0)
	for i := 0; i < len(slot); i++ {
		slotSum += slot[i]
	}

	// 内存给多少用多少，但是不超过定义的memUsage
	memUsage := alg.MemUsage
	if memUsage > mem {
		memUsage = mem

		// 由于内存不太够，理论上会影响执行速度，因此修改slotSum
		slotSum *= float64(mem) / float64(alg.MemUsage)
	}
	alg.TotalTick -= slotSum
	return 1, memUsage
}
