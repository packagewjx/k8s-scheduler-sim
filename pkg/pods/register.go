package pods

import "github.com/packagewjx/k8s-scheduler-sim/pkg/core"

func init() {
	core.RegisterPodAlgorithmFactory(BatchPod, BatchPodFactory)
	core.RegisterPodAlgorithmFactory(SimServicePod, simServicePodFacory)
}
