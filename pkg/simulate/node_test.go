package simulate

import (
	"fmt"
	"testing"
)

func TestNodeRun(t *testing.T) {
	node := NewNode(8, 1, NewFairScheduler())
	dc := NewMockDeploymentController()
	for i := 0; i < 10; i++ {
		node.AddPod(&BatchPod{
			BasePod: &BasePod{
				name:                 fmt.Sprintf("Batch-%d", i),
				podType:              "Batch",
				priority:             0,
				deploymentController: dc,
				cpuLimit:             1,
				memLimit:             1,
			},
			memUsage:  0.05,
			totalTick: 5000,
		})
	}
	for i := 0; i < 10000; i++ {
		metric := node.Tick()
		fmt.Printf("Tick: %d. CPU: %.2f%%. Mem: %.2f%%, Load: %.2f\n", i, metric.CpuUsage*100, metric.MemUsage*100, metric.Load)
	}
}
