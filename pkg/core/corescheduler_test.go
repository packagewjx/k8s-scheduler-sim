package core

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"math"
	"testing"
)

type mockAlgorithm struct {
}

func (_ *mockAlgorithm) ResourceRequest() (cpu float64, mem int) {
	return 1, 1
}

func (_ *mockAlgorithm) Tick(slot []float64, mem int) (Load float64, MemUsage int) {
	return 1, 100
}

func TestFairScheduler(t *testing.T) {
	sched, _ := GetCoreScheduler(FairScheduler)
	readyPods := make([]*Pod, 10)
	for i := 0; i < len(readyPods); i++ {
		p, _ := BuildPod(fmt.Sprintf("pod-%d", i), 1, 1, BatchPod, "null", "", v1.DefaultSchedulerName)
		pod := &Pod{
			Pod:       *p,
			CpuLimit:  1,
			MemLimit:  1,
			Algorithm: &mockAlgorithm{},
		}
		readyPods[i] = pod
	}

	cases := []struct {
		cpuState     [][]*RunEntity
		expectedSlot []float64
		expectedLen  []int
	}{
		{
			cpuState:     make([][]*RunEntity, 10),
			expectedSlot: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expectedLen:  []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			cpuState:     make([][]*RunEntity, 1),
			expectedSlot: []float64{0.1},
			expectedLen:  []int{10},
		},
		{
			cpuState:     [][]*RunEntity{},
			expectedSlot: []float64{},
			expectedLen:  []int{},
		},
		{
			cpuState:     make([][]*RunEntity, 4),
			expectedSlot: []float64{0.333333333333, 0.33333333333, 0.5, 0.5},
			expectedLen:  []int{3, 3, 2, 2},
		},
	}

	for i := 0; i < len(cases); i++ {
		newState := sched.Schedule(readyPods, cases[i].cpuState)
		if len(newState) != len(cases[i].cpuState) {
			t.Errorf("CPUState长度应该为%d", len(cases[i].cpuState))
		}
		for j := 0; j < len(newState); j++ {
			if len(newState[j]) != cases[i].expectedLen[j] {
				t.Errorf("CPU队列长度应该为%d", cases[j].expectedLen)
			}
			if math.Abs(cases[i].expectedSlot[j]-newState[j][0].Slot) > 0.00001 {
				t.Errorf("时间片大小应该为%f", cases[i].expectedSlot)
			}
		}
	}

}
