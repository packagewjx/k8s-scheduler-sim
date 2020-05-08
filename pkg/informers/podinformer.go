package informers

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type podInformer struct {
	client  kubernetes.Interface
	factory informers.SharedInformerFactory
}

func (p *podInformer) Get(name string) (*apicorev1.Pod, error) {
	return p.client.CoreV1().Pods(pkg.DefaultNamespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (p *podInformer) List(selector labels.Selector) (ret []*apicorev1.Pod, err error) {
	podList, err := p.client.CoreV1().Pods(pkg.DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}

	ret = make([]*apicorev1.Pod, 0, len(podList.Items))
	for _, pod := range podList.Items {
		if selector.Matches(labels.Set(pod.Labels)) {
			ret = append(ret, &pod)
		}
	}
	return
}

func (p *podInformer) Pods(namespace string) listerv1.PodNamespaceLister {
	return p
}

var podKeyFunc cache.KeyFunc = func(obj interface{}) (string, error) {
	return obj.(*apicorev1.Pod).Name, nil
}

func newPodInformer(client kubernetes.Interface, factory informers.SharedInformerFactory) v1.PodInformer {
	return &podInformer{
		client:  client,
		factory: factory,
	}
}

func (p *podInformer) defaultInformer(client kubernetes.Interface, _ time.Duration) cache.SharedIndexInformer {
	watcher, err := client.CoreV1().Pods(pkg.DefaultNamespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	informer, err := NewSharedIndexInformer(watcher, podKeyFunc)
	if err != nil {
		panic(err)
	}

	return informer
}

func (p *podInformer) Informer() cache.SharedIndexInformer {
	return p.factory.InformerFor(&apicorev1.Pod{}, p.defaultInformer)
}

func (p *podInformer) Lister() listerv1.PodLister {
	return p
}
