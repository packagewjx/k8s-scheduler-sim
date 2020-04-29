package simulate

import (
	"fmt"
	"testing"
)

func TestNodeRun(t *testing.T) {
	builder := NodeBuilder{
		Name:        "TestNode",
		Labels:      nil,
		Annotations: nil,
		Taints:      nil,
		CoreCount:   8,
		MemorySize:  1000,
		Scheduler:   NewFairScheduler(),
	}
	node := builder.Build()
	dc := NewMockDeploymentController()
	podBuilder := &PodBuilder{
		CpuLimit:         1,
		MemLimit:         100,
		Controller:       dc,
		AlgorithmFactory: FactoryMethod(1, 1000),
	}
	for i := 0; i < 10; i++ {
		podBuilder.Name = fmt.Sprintf("Batch-%d", i)
		node.AddPod(podBuilder.Build())
	}
	for i := 0; i < 10000; i++ {
		metric := node.Tick()
		fmt.Printf("Tick: %d. CPU: %.2f%%. Mem: %.2f%%, Load: %.2f\n", i, metric.CpuUsage*100, metric.MemUsage*100, metric.Load)
	}
}
