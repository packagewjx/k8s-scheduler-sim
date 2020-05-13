package core

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/metrics"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

const (
	NodeAnnotationCoreScheduler = "github.com/packagewjx/corealgorithm"
)

// Node 集群中的模拟Node。继承Kubernetes的Node，但是字段仅保留与调度有关的信息，也就是Kubernetes调度器使用到的部分。
type Node struct {
	v1.Node
	Scheduler CoreScheduler
	// 本节点上的Pod
	Pods map[string]*Pod
	// CpuState 反映当前CPU状态，每个CPU上有自己的RunEntity队列，表示在上一个周期中运行的所有Pod
	CpuState [][]*RunEntity

	// 上一轮的CPU使用百分比，用于查看是否有资源竞争，模拟高资源竞争时CPU处理能力的下降
	LastCpuUsage float64

	bindLock sync.Mutex
}

func BuildNode(name string, cpu, mem, numPods, coreScheduler string) *v1.Node {
	return &v1.Node{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				NodeAnnotationCoreScheduler: coreScheduler,
			},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Capacity: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(cpu),
				v1.ResourceMemory: resource.MustParse(mem),
				v1.ResourcePods:   resource.MustParse(numPods),
			},
			Allocatable: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(cpu),
				v1.ResourceMemory: resource.MustParse(mem),
				v1.ResourcePods:   resource.MustParse(numPods),
			},
		},
	}
}

// TODO 可能会有一些绑定失败的条件，如内存不够用等
func (n *Node) BindPod(pod *Pod) error {
	n.bindLock.Lock()
	defer n.bindLock.Unlock()

	logrus.Infof("Binding Pod %s to Node %s", pod.Name, n.Name)

	// 暂时使用pod.Name作为键
	key, _ := PodKeyFunc(pod)
	if _, ok := n.Pods[key]; ok {
		return fmt.Errorf("pod %s already bounded", key)
	}
	n.Pods[key] = pod
	return nil
}

func (n *Node) EvictPod(pod *Pod) error {
	delete(n.Pods, pod.Name)
	return nil
}

// 根据节点拥有的Pod，更新当前的节点状态，包括资源使用率，Pod状态等
// TODO 内存分配
func (n *Node) Tick(client kubernetes.Interface) *metrics.TickMetrics {
	terminatedPods := make([]*Pod, 0)
	type PodResource struct {
		slot     []float64
		mem      int
		load     float64
		memUsage int
	}

	readyPods := make([]*Pod, 0, len(n.Pods))
	podResource := make([]*PodResource, 0, len(n.Pods))
	podIdxMap := make(map[*Pod]int)

	// 查看Pod的状态
	for _, pod := range n.Pods {
		if pod.Status.Phase == v1.PodRunning {
			_, mem := pod.Algorithm.ResourceRequest()
			podIdxMap[pod] = len(podResource)
			readyPods = append(readyPods, pod)
			podResource = append(podResource, &PodResource{
				slot: make([]float64, 0),
				mem:  mem,
			})
		} else if pod.Status.Phase == v1.PodPending {
			pod.Status.Phase = v1.PodRunning
			pod.Spec.NodeName = n.Name
			_, err := client.CoreV1().Pods(DefaultNamespace).UpdateStatus(context.TODO(), &pod.Pod, metav1.UpdateOptions{})
			if err != nil {
				logrus.Errorf("Error Updating Pod %s Status", pod.Name)
			}
		} else {
			terminatedPods = append(terminatedPods, pod)
		}
	}

	logrus.Debugf("Node %s Terminating Pods", n.Name)
	for _, pod := range terminatedPods {
		logrus.Infof("Pod %s is now %s", pod.Name, pod.Status.Phase)
		// TODO 通过watch.Interface通知Pods结束
		_, err := client.CoreV1().Pods(DefaultNamespace).UpdateStatus(context.TODO(), &pod.Pod, metav1.UpdateOptions{})
		if err != nil {
			logrus.Errorf("Node %s Update pod status for pod %s error: %v", n.Name, pod.Name, err)
		}
		// 从本节点移除
		logrus.Tracef("Removing Pod %s from Node %s", pod.Name, n.Name)
		key, _ := PodKeyFunc(pod)
		delete(n.Pods, key)
	}

	// 执行调度算法
	logrus.Debugf("Node %s Calculating cpu slots for ready pods", n.Name)
	cpuState := n.Scheduler.Schedule(readyPods, n.CpuState)
	for i := 0; i < len(cpuState); i++ {
		for j := 0; j < len(cpuState[i]); j++ {
			entity := cpuState[i][j]
			podIdx := podIdxMap[entity.Pod]
			podResource[podIdx].slot = append(podResource[podIdx].slot, entity.Slot)
		}
	}

	// 根据分配结果更新Pod的执行状态
	logrus.Debugf("Node %s Updating Pod status", n.Name)
	memUsed := 0
	load := float64(0)
	for i := 0; i < len(readyPods); i++ {
		logrus.Tracef("Node %s Updating Pod %s status", n.Name, readyPods[i].Name)
		cpuPressureReduction(podResource[i].slot, n.LastCpuUsage)
		podResource[i].load, podResource[i].memUsage = readyPods[i].Algorithm.Tick(podResource[i].slot, podResource[i].mem)

		// 计算统计
		memUsed += podResource[i].memUsage
		load += podResource[i].load
	}

	// 更新节点的状态

	// 根据Pod返回的负载信息，计算CPU统计数据
	logrus.Debugf("Updating Node %s cpu and mem usage", n.Name)
	cpuUsed := float64(0)
	coreCount, _ := n.Status.Capacity.Cpu().AsInt64()
	memSize, _ := n.Status.Capacity.Memory().AsInt64()
	for i := 0; i < int(coreCount); i++ {
		for j := 0; j < len(cpuState[i]); j++ {
			entity := cpuState[i][j]
			podIdx := podIdxMap[entity.Pod]
			cpuUsed += entity.Slot * podResource[podIdx].load
		}
	}

	cpuUsage := cpuUsed / float64(coreCount)
	n.LastCpuUsage = cpuUsage

	n.Status.Allocatable.Cpu().Set(coreCount - int64(cpuUsed))
	n.Status.Allocatable.Memory().Set(memSize - int64(memUsed))
	_, err := client.CoreV1().Nodes().UpdateStatus(context.TODO(), &n.Node, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Update Node %s Status error: %v", n.Name, err)
	}

	return &metrics.TickMetrics{
		CpuUsage: cpuUsage,
		MemUsage: float64(memUsed) / float64(memSize),
		Load:     load,
	}

}

// 当高CPU压力时，减少实际slot数，反映CPU速率的下降。注意在最后统计时，使用没有减少的slot。
func cpuPressureReduction(slot []float64, cpuUsage float64) {
	// 只在大于0.7时触发
	if cpuUsage >= 0.7 {
		// 线性变化，减少时间片
		reduction := 1 - (cpuUsage - 0.7)

		for i := 0; i < len(slot); i++ {
			slot[i] *= reduction
		}
	}
}
