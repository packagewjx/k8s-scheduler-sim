package simulate

import (
	"fmt"
)

// mockDeploymentController 虚拟的控制器，用于测试
type mockDeploymentController struct {
	tickNum int
}

func NewMockDeploymentController() DeploymentController {
	return &mockDeploymentController{0}
}

func (ctrl *mockDeploymentController) Tick() (addPod []*Pod, removePod []*Pod) {
	panic("implement me")
}

func (ctrl *mockDeploymentController) InformPodEvent(event *PodEvent) {
	fmt.Printf("%s has terminated.\n", event.Who.Name)
}
