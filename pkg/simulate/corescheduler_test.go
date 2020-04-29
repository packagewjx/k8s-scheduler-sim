package simulate

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

type mockPod struct {
	name string
}

func (pod *mockPod) Name() string {
	if pod.name == "" {
		pod.name = fmt.Sprintf("Mock-%f", rand.Float64())
	}
	return pod.name
}

func (pod *mockPod) Priority() int {
	return 0
}

func (pod *mockPod) Type() string {
	return "Mock"
}

func (pod *mockPod) ResourceLimit() (cpuLimit int, memLimit float64) {
	return 1, 1
}

func (pod *mockPod) ResourceRequest() (cpu int, mem float64) {
	return 1, 0.01
}

func (pod *mockPod) Tick(_ []float64, _ float64) (Load, MemUsage float64) {
	return 1, 0.01
}

func (pod *mockPod) GetState() PodState {
	return RunningState
}

func (pod *mockPod) DeploymentController() DeploymentController {
	return nil
}

func TestFairScheduler(t *testing.T) {
	sched := NewFairScheduler()
	readyPods := make([]Pod, 10)
	for i := 0; i < len(readyPods); i++ {
		readyPods[i] = &mockPod{}
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
