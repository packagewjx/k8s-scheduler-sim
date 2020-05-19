package pods

import (
	"encoding/json"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"k8s.io/api/core/v1"
	"math"
	"testing"
)

func TestServicePod(t *testing.T) {
	baseMem := int64(300 * (1 << 20))
	argByte, _ := json.Marshal(&SimServicePodArgs{BaseMem: baseMem})
	pod := &core.Pod{
		Pod:       v1.Pod{},
		CpuLimit:  4,
		MemLimit:  1 << 30,
		Algorithm: nil,
	}
	alg, _ := simServicePodFacory(string(argByte), pod)
	pod.Algorithm = alg
	servAlg := alg.(ServicePod)

	slot := []float64{1, 1, 1, 1}

	// 初始化
	load, memUsage := servAlg.Tick(slot, 1<<30)
	if load > 0.5 {
		t.Error("load should be 0.5")
	}
	if memUsage != baseMem {
		t.Error("should be baseMem")
	}

	service := &ServiceContext{
		RequestId:    0,
		OnDone:       func(requestId int) {},
		SlotRequired: 0.01,
		MemRequired:  500 * (1 << 10),
	}
	//开始服务
	for i := 0; i < 1000; i++ {
		err := servAlg.DeliverRequest(service)
		if err != nil {
			t.Error("should not be:", err)
		}
	}
	for load, memUsage := servAlg.Tick(slot, 1<<30); load != 0; load, memUsage = servAlg.Tick(slot, 1<<30) {
		fmt.Println(load, memUsage)
	}

}

func TestSingleCPU(t *testing.T) {
	argByte, _ := json.Marshal(&SimServicePodArgs{BaseMem: 0})
	pod := &core.Pod{
		Pod:       v1.Pod{},
		CpuLimit:  1,
		MemLimit:  1 << 30,
		Algorithm: nil,
	}
	alg, _ := simServicePodFacory(string(argByte), pod)
	pod.Algorithm = alg
	servAlg := alg.(ServicePod)
	servAlg.Tick([]float64{1}, 1<<30)

	// 单CPU测试
	memCases2 := 4504 // 由于精度丢失，不能给表达式，这里直接列出结果
	testCases := []struct {
		doneId    []int
		services  []*ServiceContext
		tickCount int
		tickLoad  []float64
		tickMem   []int64
	}{
		{
			doneId: []int{0, 1, 2, 3, 4},
			services: []*ServiceContext{
				{
					RequestId:    0,
					OnDone:       nil,
					SlotRequired: 0.1,
					MemRequired:  1 << 10,
				},
				{
					RequestId:    1,
					OnDone:       nil,
					SlotRequired: 0.1,
					MemRequired:  1 << 10,
				},
				{
					RequestId:    2,
					OnDone:       nil,
					SlotRequired: 0.1,
					MemRequired:  1 << 10,
				},
				{
					RequestId:    3,
					OnDone:       nil,
					SlotRequired: 0.1,
					MemRequired:  1 << 10,
				},
				{
					RequestId:    4,
					OnDone:       nil,
					SlotRequired: 0.1,
					MemRequired:  1 << 10,
				},
			},
			tickCount: 1,
			tickLoad:  []float64{0.5},
			tickMem:   []int64{1 << 10},
		},
		{
			doneId: []int{0, 1, 2, 3},
			services: []*ServiceContext{
				{
					RequestId:    0,
					OnDone:       nil,
					SlotRequired: 0.4,
					MemRequired:  1 << 10,
				},
				{
					RequestId:    1,
					OnDone:       nil,
					SlotRequired: 0.4,
					MemRequired:  1 << 12,
				},
				{
					RequestId:    2,
					OnDone:       nil,
					SlotRequired: 0.4,
					MemRequired:  1 << 10,
				},
				{
					RequestId:    3,
					OnDone:       nil,
					SlotRequired: 0.5,
					MemRequired:  1 << 10,
				},
			},
			tickCount: 2,
			tickLoad:  []float64{0.8, 0.9},
			tickMem:   []int64{int64(memCases2), 1 << 10},
		},
	}

	for _, testCase := range testCases {
		shouldBe := 0
		onDone := func(requestId int) {
			if requestId != testCase.doneId[shouldBe] {
				t.Errorf("done id should be %d not %d", testCase.doneId[shouldBe], requestId)
			}
			shouldBe++
		}
		for _, service := range testCase.services {
			service.OnDone = onDone
			err := servAlg.DeliverRequest(service)
			if err != nil {
				t.Error("Should not be error", err)
			}
		}

		for tick := 0; tick < testCase.tickCount; tick++ {
			load, mem := servAlg.Tick([]float64{1}, 1<<30)
			if !floatEquals(load, testCase.tickLoad[tick]) {
				t.Errorf("load should be %f not %f", testCase.tickLoad[tick], load)
			}
			if mem != testCase.tickMem[tick] {
				t.Errorf("mem should be %d not %d", testCase.tickMem[tick], mem)
			}
		}

		load, _ := servAlg.Tick([]float64{1}, 1<<30)
		if load != 0 {
			t.Errorf("Tick error")
		}
	}
}

func TestMaxProbing(t *testing.T) {
	argByte, _ := json.Marshal(&SimServicePodArgs{BaseMem: 0})
	pod := &core.Pod{
		Pod:       v1.Pod{},
		CpuLimit:  1,
		MemLimit:  1 << 30,
		Algorithm: nil,
	}
	alg, _ := simServicePodFacory(string(argByte), pod)
	pod.Algorithm = alg
	servAlg := alg.(ServicePod)
	servAlg.Tick([]float64{1}, 1<<30)

	// Case 1
	onDone := func(requestId int) {
		if requestId != 1 {
			t.Error("Id 不对")
		}
	}
	for i := 0; i < 9; i++ {
		_ = servAlg.DeliverRequest(&ServiceContext{
			RequestId:    0,
			OnDone:       onDone,
			SlotRequired: 0.5,
			MemRequired:  0,
		})
	}
	_ = servAlg.DeliverRequest(&ServiceContext{
		RequestId:    1,
		OnDone:       onDone,
		SlotRequired: 0.2,
		MemRequired:  0,
	})
	load, _ := servAlg.Tick([]float64{0.3}, 1<<10)
	if !floatEquals(0.2/0.3, load) {
		t.Errorf("wrong load")
	}

}

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}
