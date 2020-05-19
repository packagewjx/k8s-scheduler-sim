package core

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/cache"
	"testing"
	"time"
)

func newFakePod(name string) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				PodAnnotationCpuLimit:             "1",
				PodAnnotationMemLimit:             "1",
				PodAnnotationDeploymentController: "null",
				PodAnnotationAlgorithm:            testPod,
				PodAnnotationInitialState:         "",
			},
			UID: uuid.NewUUID(),
		},
		Spec: v1.PodSpec{SchedulerName: v1.DefaultSchedulerName},
		Status: v1.PodStatus{
			Phase: v1.PodPending,
		},
	}
}

func TestScheduleOne(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	sim := NewSchedulerSimulator(1000)

	node := &v1.Node{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "node-1",
			ClusterName: "testcluster",
			Annotations: map[string]string{
				NodeAnnotationCoreScheduler: FairScheduler,
			},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Capacity: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("8"),
				v1.ResourceMemory: resource.MustParse("16G"),
				v1.ResourcePods:   resource.MustParse("100"),
			},
			Allocatable: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("8"),
				v1.ResourceMemory: resource.MustParse("16G"),
				v1.ResourcePods:   resource.MustParse("100"),
			},
			Phase: v1.NodeRunning,
		},
	}
	node, err := sim.GetKubernetesClient().CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("node create fail: %v", err)
	}

	// wait for node creation event
	<-time.After(50 * time.Millisecond)

	podName := "fakepod"
	pod := newFakePod(podName)

	pod, err = sim.GetKubernetesClient().CoreV1().Pods(DefaultNamespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("create pod failed: %v", err)
	}

	select {
	case str := <-sim.(*schedSim).bindPodCh:
		if str[0] != 'T' {
			t.Error("Schedule failed")
		}
	case <-time.After(time.Second):
		t.Error("Schedule failed")
	}
}

func TestNodeClient(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	sim := NewSchedulerSimulator(1000)
	defer sim.(*schedSim).cancelFunc()

	nodeClient := sim.GetKubernetesClient().CoreV1().Nodes()

	informer := sim.GetInformerFactory().Core().V1().Nodes().Informer()

	ch := make(chan bool)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			t.Log("added")
			ch <- true
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			t.Log("udpated")
			ch <- true
		},
		DeleteFunc: func(obj interface{}) {
			t.Log("deleted")
			ch <- true
		},
	})

	testNodeName := "test-1"
	node := &v1.Node{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:            testNodeName,
			ResourceVersion: "1",
			ClusterName:     "TestCluster",
			Annotations: map[string]string{
				NodeAnnotationCoreScheduler: FairScheduler,
			},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Capacity: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("4"),
				v1.ResourceMemory: resource.MustParse("4G"),
			},
			Allocatable: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("4"),
				v1.ResourceMemory: resource.MustParse("4G"),
			},
			Phase: v1.NodeRunning,
		},
	}

	// test create function
	node, err := nodeClient.Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal("node is nil")
	}

	if node.Name != testNodeName {
		t.Errorf("name is %s not %s", node.Name, testNodeName)
	}

	// check add event
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("no add event")
	}

	// test get function
	node, err = nodeClient.Get(context.TODO(), testNodeName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if node.Name != testNodeName {
		t.Fatalf("name is %s not %s", node.Name, testNodeName)
	}

	// test list function
	list, err := nodeClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Errorf("list error: %v", err)
	}
	if len(list.Items) != 1 {
		t.Errorf("list length not 1")
	}
	if list.Items[0].Name != testNodeName {
		t.Errorf("list node name incorrect")
	}

	// test update function
	node.Labels = map[string]string{
		"testing": "true",
	}
	_, err = nodeClient.Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		t.Error(err)
	}

	node, err = nodeClient.Get(context.TODO(), testNodeName, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	if node.Labels["testing"] != "true" {
		t.Error("update is not success")
	}

	// check update event
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("no update event")
	}

	// test updatestatus function
	node.Status.Allocatable[v1.ResourceMemory] = resource.MustParse("2G")
	_, err = nodeClient.UpdateStatus(context.TODO(), node, metav1.UpdateOptions{})
	node, err = nodeClient.Get(context.TODO(), testNodeName, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	if node.Status.Allocatable.Memory().Cmp(resource.MustParse("2G")) != 0 {
		t.Error("update status is not success")
	}

	// check update event
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("no update event")
	}

	// test delete function
	err = nodeClient.Delete(context.TODO(), testNodeName, metav1.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
	node, err = nodeClient.Get(context.TODO(), testNodeName, metav1.GetOptions{})
	if node != nil {
		t.Error("should be nil")
	}

	// check delete event
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("no delete event")
	}

}

func TestPodClient(t *testing.T) {
	sim := NewSchedulerSimulator(1000)
	defer sim.(*schedSim).cancelFunc() // inorder to stop InformerFactory and scheduler

	podClient := sim.GetKubernetesClient().CoreV1().Pods(DefaultNamespace)
	podInformer := sim.GetInformerFactory().Core().V1().Pods().Informer()
	ch := make(chan bool)
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			t.Log("added")
			ch <- true
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			t.Log("updated")
			ch <- true
		},
		DeleteFunc: func(obj interface{}) {
			t.Log("deleted")
			ch <- true
		},
	})

	podName := "test-1"
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: map[string]string{},
			Annotations: map[string]string{
				PodAnnotationDeploymentController: "null",
				PodAnnotationAlgorithm:            testPod,
				PodAnnotationMemLimit:             "1000",
				PodAnnotationCpuLimit:             "1",
			},
		},
		Spec: v1.PodSpec{},
		Status: v1.PodStatus{
			Phase: v1.PodPending,
		},
	}

	// test create
	pod, err := podClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if pod == nil {
		t.Fatalf("created pod is nil")
	}

	select {
	case <-ch:
	case <-time.After(10 * time.Millisecond):
		t.Error("no add event")
	}

	// test get
	pod, err = podClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get pod failed: %v", err)
	}
	if pod == nil {
		t.Fatalf("get pod nil")
	}

	// test list
	list, err := podClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Errorf("list error: %v", err)
	}
	if len(list.Items) != 1 {
		t.Errorf("not 1 item")
	}
	if list.Items[0].Name != podName {
		t.Errorf("list pod name incorrect")
	}

	// test update
	pod.Labels["testing"] = "true"
	_, err = podClient.Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	pod, _ = podClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if pod.Labels["testing"] != "true" {
		t.Errorf("update failed")
	}

	select {
	case <-ch:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("no update event")
	}

	// test updatestatus
	pod.Status.Phase = v1.PodFailed
	_, err = podClient.UpdateStatus(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		t.Errorf("pod updatestatus error: %v", err)
	}
	pod, _ = podClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if pod.Status.Phase != v1.PodFailed {
		t.Errorf("pod updatestatus failed")
	}

	select {
	case <-ch:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("no update event")
	}

	// test delete
	err = podClient.Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		t.Errorf("delete failed: %v", err)
	}
	pod, err = podClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if err == nil {
		t.Errorf("should be error")
	}
	if pod != nil {
		t.Errorf("should be nil")
	}

	select {
	case <-ch:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("no delete event")
	}
}

func TestDeploy10Tick(t *testing.T) {
	simulator := NewSchedulerSimulator(200)
	node := BuildNode("node-1", "10", "2G", "5000", FairScheduler)
	_, err := simulator.GetKubernetesClient().CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	simulator.RegisterBeforeUpdateController(&deploy10TimesController{
		phase: 0,
		sim:   simulator,
	})

	simulator.Run()
}

func TestDeployMultiplePods(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	sim := NewSchedulerSimulator(300)
	node := BuildNode("node-1", "10", "100G", "1000", FairScheduler)
	_, err := sim.GetKubernetesClient().CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	sim.RegisterBeforeUpdateController(&deployMultiplePodsController{
		sim:    sim,
		podNum: 1000,
	})

	sim.Run()
}

func TestDeployPodExceedLimit(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	sim := NewSchedulerSimulator(1000)
	node := BuildNode("node-1", "10", "100G", "10", FairScheduler)
	_, err := sim.GetKubernetesClient().CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	sim.RegisterBeforeUpdateController(&deployMultiplePodsController{
		sim:    sim,
		podNum: 100,
	})

	sim.Run()
}

type deploy10TimesController struct {
	phase int
	sim   SchedulerSimulator
}

func (m *deploy10TimesController) Tick() {
	if m.phase < 10 {
		pod := newFakePod(fmt.Sprintf("pod-%d", m.phase))
		_, err := m.sim.GetKubernetesClient().CoreV1().Pods(DefaultNamespace).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		m.phase++
	} else {
		m.sim.DeleteBeforeController(m)
	}
}

type deployMultiplePodsController struct {
	sim    SchedulerSimulator
	podNum int
}

func (d *deployMultiplePodsController) Tick() {
	for i := 0; i < d.podNum; i++ {
		pod := newFakePod(fmt.Sprintf("pod-%d", i))
		_, err := d.sim.GetKubernetesClient().CoreV1().Pods(DefaultNamespace).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
	}
	d.sim.DeleteBeforeController(d)
}
