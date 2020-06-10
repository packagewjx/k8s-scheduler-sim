package controllers

import (
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/core"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/pods"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/cache"
	"math"
)

const LabelService = "service"

// ServiceController 模拟Kubernetes中进行负载均衡的Service，负责接收外部请求，并分配到各个Service Pod上。
type ServiceController interface {
	core.Controller
}

// ServiceContextFactory 负责在每一次的Tick的时候返回需要分发的服务
type ServiceContextFactory func() []*pods.ServiceContext

// NewServiceController podTemplate需要正确设置Annotation的CPU与内存，其余Annotation将由此Controller设置
func NewServiceController(sim core.SchedulerSimulator, name string, podNum int, factory ServiceContextFactory, podTemplate *v1.Pod) ServiceController {
	if podTemplate.Annotations == nil {
		podTemplate.Annotations = make(map[string]string)
	}
	podTemplate.Annotations[core.PodAnnotationAlgorithm] = pods.SimServicePod
	podTemplate.Annotations[core.PodAnnotationInitialState] = ""
	podTemplate.Annotations[core.PodAnnotationDeploymentController] = name
	if podTemplate.Labels == nil {
		podTemplate.Labels = make(map[string]string)
	}
	// 标记这个Pod是本Service构建的
	podTemplate.Labels[LabelService] = name

	rc := NewReplicationController(sim, fmt.Sprintf("replication-controller-%s", name), podNum, func() *v1.Pod {
		uid := uuid.NewUUID()
		pod := podTemplate.DeepCopy()

		pod.Name = fmt.Sprintf("service_pod_%s_%s", name, uid)
		pod.UID = uid
		return pod
	})
	sim.RegisterBeforeUpdateController(rc)

	return &serviceController{
		podTemplate: podTemplate,
		name:        name,
		sim:         sim,
		factory:     factory,
		pods:        make(map[string]*core.Pod),
		rc:          rc,
		initialized: false,
		tick:        0,
		requestTick: make(map[int]int),
		queue:       make([]*pods.ServiceContext, 0, 10),
		met:         &serviceMetrics{},
	}
}

type serviceController struct {
	podTemplate *v1.Pod
	name        string
	sim         core.SchedulerSimulator
	factory     ServiceContextFactory
	pods        map[string]*core.Pod
	// rc 用于维护Pod数量
	rc          ReplicationController
	initialized bool
	// tick 记录当前的tick数
	tick int
	// requestTick 记录请求何时到来，用于统计。键为RequestId，值为tick
	requestTick map[int]int
	// queue 未处理的队列
	queue []*pods.ServiceContext
	// met 用于统计
	met *serviceMetrics
}

func (c *serviceController) Name() string {
	return c.name
}

func (c *serviceController) Tick() {
	if !c.initialized {
		// 注册Informer，加入新的Pod，删除停止的Pod与重新部署新的Pod
		c.sim.GetInformerFactory().Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: nil,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*v1.Pod)
				newPod := newObj.(*v1.Pod)
				// 有没有可能Label被改，应该不会的
				if oldPod.Labels == nil || oldPod.Labels[LabelService] != c.name {
					return
				}
				// 调度成功后再加入到可赋值集群
				if oldPod.Spec.NodeName == "" && isPodBindSuccess(newPod) {
					logrus.Infof("Service %s: Pod %s is ready for handling requests.", c.name, newPod.Name)
					pod, err := c.sim.GetPod(oldPod.Name)
					if err != nil {
						logrus.Errorf("Service %s: Error getting pod %s in create event handler: %v", c.name, oldPod.Name, err)
						return
					}
					c.pods[pod.Name] = pod
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if pod.Labels != nil && pod.Labels[LabelService] == c.name {
					delete(c.pods, pod.Name)
				}
			},
		})
		c.initialized = true
		return
	}

	requests := c.factory()
	for _, req := range requests {
		req.OnDone = c.onDone
		c.requestTick[req.RequestId] = c.tick
		c.queue = append(c.queue, req)
	}
	requests = nil

	// 仅当有Pod时分发
	if len(c.pods) > 0 {
		sumRequests := 0
		for _, pod := range c.pods {
			sumRequests += pod.Algorithm.(pods.ServicePod).GetRequestQueueLen()
		}
		avgRequest := (sumRequests + len(c.queue)) / len(c.pods)
		// 记录c.queue的位置，应该不会超过len(c.queue)
		idx := 0
		for _, pod := range c.pods {
			alg := pod.Algorithm.(pods.ServicePod)
			for i := alg.GetRequestQueueLen(); i < avgRequest; i++ {
				err := alg.DeliverRequest(c.queue[idx])
				if err != nil {
					logrus.Warnf("Service %s: Pod %s can not handle requests. Reason: %v", c.name, pod.Name, err)
					break
				}
				idx++
			}
		}

		// 将前面已经分发的请求删除
		c.queue = c.queue[idx:]
	}

	c.tick++
}

func (c *serviceController) onDone(requestId int) {
	tick, ok := c.requestTick[requestId]
	if ok {
		c.met.add(uint8(c.tick - tick))
	} else {
		logrus.Warnf("Service %s: No requestId %d", c.name, requestId)
	}
}

type serviceMetrics struct {
	arr []uint8
	sum int
}

func (s *serviceMetrics) add(data uint8) {
	s.arr = append(s.arr, data)
	s.sum += int(data)
}

func (s *serviceMetrics) getMetrics() (avgTicks float64, p50, p90, p99 uint8) {
	if len(s.arr) == 0 {
		return
	}

	avgTicks = float64(s.sum) / float64(len(s.arr))
	p50 = findPercentile(s.arr, 50)
	p90 = findPercentile(s.arr, 90)
	p99 = findPercentile(s.arr, 99)
	return
}

func isPodBindSuccess(pod *v1.Pod) bool {
	if pod.Spec.NodeName != "" {
		for _, condition := range pod.Status.Conditions {
			if condition.Reason == v1.PodReasonUnschedulable || (condition.Type == v1.PodScheduled && condition.Status == v1.ConditionFalse) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func findPercentile(arr []uint8, percentile float64) uint8 {
	p := int(math.Ceil(float64(len(arr)-1) / 100 * percentile))

	start := 0
	end := len(arr) - 1
	for mid := partition(arr, start+(end-start)/2, start, end); mid != p; mid = partition(arr, start+(end-start)/2, start, end) {
		if mid > p {
			end = mid - 1
		}
		if mid < p {
			start = mid + 1
		}
	}

	return arr[p]
}

func partition(arr []uint8, pivotPos, start, end int) int {
	if start >= end {
		return start
	}

	i := start
	j := end
	pivotVal := arr[pivotPos]
	arr[pivotPos] = arr[start]

	for i < j {
		for i < j && arr[j] >= pivotVal {
			j--
		}
		arr[i] = arr[j]

		for i < j && arr[i] < pivotVal {
			i++
		}
		arr[j] = arr[i]
	}

	arr[i] = pivotVal
	return i
}
