package controllers

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

func NewReplicationController(sim core.SchedulerSimulator, controllerName string, replicaNum int, podFactory func() *v1.Pod) ReplicationController {
	replicas := make(map[string]*v1.Pod)
	r := &replicationController{
		name:       controllerName,
		sim:        sim,
		replicas:   replicas,
		podFactory: podFactory,
		replicaNum: replicaNum,
		state:      initializing,
	}

	return r
}

const LabelReplicationController = "github.com/packagewjx/replicationcontroller"

type ReplicationController interface {
	core.Controller
	SetReplicaNum(num int)
	Terminate()
}

type controllerState string

var (
	initializing        controllerState = "initializing"
	running             controllerState = "running"
	terminated          controllerState = "terminated"
	terminating         controllerState = "terminating"
	waitForPodTerminate controllerState = "waitForPodTerminate"
)

type replicationController struct {
	// name 本控制器的名称，通常代表所控制的服务
	name string
	// sim 集群接口
	sim core.SchedulerSimulator
	// replicas 记录本控制器部署的Pod
	replicas map[string]*v1.Pod
	// stopping 记录哪些Pod正在停止，防止多次停止。若有元素，则代表已经执行过停止Pod的代码，在停止之前，不会继续停止其他Pod。
	stopping map[string]bool
	// replicaNum 期望的副本数量
	replicaNum int
	// podFactory 构造新的Pod的函数，本Pod不能重名，否则会有错误
	podFactory func() *v1.Pod
	// state 记录本控制器的状态
	state controllerState
}

func (r *replicationController) Name() string {
	return r.name
}

func (r *replicationController) Tick() {
	switch r.state {
	case initializing:
		// 注册监听器
		r.sim.GetInformerFactory().Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				logrus.Debugf("ReplicationController %s: Pod %s added successfully", r.name, pod.Name)
				if pod.Labels != nil && pod.Labels[LabelReplicationController] == r.name {
					r.replicas[pod.Name] = pod
				}
			},
			UpdateFunc: func(_, newObj interface{}) {
				pod := newObj.(*v1.Pod)
				logrus.Debugf("ReplicationController %s: Pod %s updated", r.name, pod.Name)
				if pod.Labels != nil && pod.Labels[LabelReplicationController] == r.name {
					r.replicas[pod.Name] = pod
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				logrus.Debugf("ReplicationController %s: Pod %s deleted", r.name, pod.Name)
				if pod.Labels != nil && pod.Labels[LabelReplicationController] == r.name {
					delete(r.replicas, pod.Name)
					if r.stopping[pod.Name] {
						delete(r.stopping, pod.Name)
					}
				}
			},
		})
		// Ensure handler has successfully created
		<-time.After(100 * time.Millisecond)
		r.state = running
	case terminated:
		// r.sim.DeleteBeforeController(r)
		return
	case terminating:
		logrus.Infof("ReplicationController %s: now entering terminating state.", r.name)
		for podName, _ := range r.replicas {
			logrus.Infof("ReplicationController %s: Terminating pod %s", r.name, podName)
			err := r.sim.GetKubernetesClient().CoreV1().Pods(core.DefaultNamespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
			if err != nil {
				logrus.Errorf("ReplicationController %s: Error deleting pods %s: %v", r.name, podName, err)
			}
		}
		r.state = waitForPodTerminate
	case waitForPodTerminate:
		logrus.Infof("ReplicationController %s: waiting for pod termination.", r.name)
		terminatedPod := make([]*v1.Pod, 0, 10)
		for _, pod := range r.replicas {
			if pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodSucceeded {
				logrus.Infof("ReplicationController %s: Pod %s terminated", r.name, pod.Name)
				terminatedPod = append(terminatedPod, pod)
			}
		}
		for _, pod := range terminatedPod {
			err := r.sim.GetKubernetesClient().CoreV1().Pods(core.DefaultNamespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				logrus.Errorf("ReplicationController %s: delete pod %s failed: %v", r.name, pod.Name, err)
				// delete next tick
				continue
			}
			delete(r.replicas, pod.Name)
		}

		if len(r.replicas) == 0 {
			logrus.Infof("ReplicationController %s terminate successfully", r.name)
			r.state = terminated
		}
	case running:
		terminated := make([]*v1.Pod, 0, 10)
		// 首先去除已经Terminated的
		for _, pod := range r.replicas {
			if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
				logrus.Infof("ReplicationController %s: Pod %s terminated", r.name, pod.Name)
				terminated = append(terminated, pod)
			}
		}
		for _, pod := range terminated {
			delete(r.replicas, pod.Name)
		}

		if len(r.replicas) > r.replicaNum {
			if len(r.stopping) == 0 {
				logrus.Infof("ReplicationController %s: Pod replica is larger than expected replica number, terminating pods.", r.name)
				// 若stopping没有Pod时才停止，否则等待。随机选择几个停止。
				for _, pod := range r.replicas {
					logrus.Infof("ReplicationController %s: Terminating pod %s", r.name, pod.Name)
					err := r.sim.GetKubernetesClient().CoreV1().Pods(core.DefaultNamespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
					if err != nil {
						logrus.Errorf("ReplicationController %s: error deleting pod %s: %v", r.name, pod.Name, err)
					}
					r.stopping[pod.Name] = true
				}
			}
		} else if len(r.replicas) < r.replicaNum {
			logrus.Infof("ReplicationController %s: Pod is not enough, creating pods.", r.name)
			// 启动Pod，以满足数量
			for len(r.replicas) < r.replicaNum {
				pod := r.podFactory()
				// 添加本Controller的标志，以确定时本Controller部署的
				if pod.Labels == nil {
					pod.Labels = make(map[string]string)
				}
				pod.Labels[LabelReplicationController] = r.name
				podName := pod.Name

				logrus.Infof("ReplicationController %s: Creating pod %s", r.name, pod.Name)
				pod, err := r.sim.GetKubernetesClient().CoreV1().Pods(core.DefaultNamespace).Create(context.TODO(), pod, metav1.CreateOptions{})
				if err != nil {
					logrus.Errorf("ReplicationController %s: error creating pod %s: %v", r.name, podName, err)
					continue
				}
				logrus.Debugf("ReplicationController %s: Pod %s created", r.name, pod.Name)
				r.replicas[pod.Name] = pod
			}
		}
	}
}

func (r *replicationController) SetReplicaNum(num int) {
	r.replicaNum = num
}

func (r *replicationController) Terminate() {
	r.state = terminating
}
