package pods

import (
	"encoding/json"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

// batchPodAlgorithm 模拟批处理任务的Pod，默认一直在跑任务，因此负载一直为1，内存使用基本固定
type batchPodAlgorithm struct {
	Pod           *core.Pod
	markTerminate bool
	MemUsage      int64
	TotalTick     float64
}

func (alg *batchPodAlgorithm) Terminate() {
	alg.markTerminate = true
}

type BatchPodState struct {
	MemUsage  int64   `json:"memUsage"`
	TotalTick float64 `json:"totalTick"`
}

const BatchPod = "BatchPod"

var BatchPodFactory core.PodAlgorithmFactory = func(stateJson string, pod *core.Pod) (core.PodAlgorithm, error) {
	state := &BatchPodState{}
	if stateJson != "" {
		err := json.Unmarshal([]byte(stateJson), state)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing state json")
		}
	}

	return &batchPodAlgorithm{
		Pod:           pod,
		MemUsage:      state.MemUsage,
		TotalTick:     state.TotalTick,
		markTerminate: false,
	}, nil
}

func (alg *batchPodAlgorithm) ResourceRequest() (cpu float64, mem int64) {
	return alg.Pod.CpuLimit, alg.MemUsage
}

func (alg *batchPodAlgorithm) Tick(slot []float64, mem int64) (Load float64, MemUsage int64) {
	if alg.TotalTick < 0 || alg.markTerminate {
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
