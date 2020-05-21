package informers

import (
	"context"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
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

type nodeInformer struct {
	client       kubernetes.Interface
	factory      informers.SharedInformerFactory
	listenerChan chan cache.ResourceEventHandler
	store        cache.Store
}

func NewNodeInformer(client kubernetes.Interface, factory informers.SharedInformerFactory) v1.NodeInformer {
	return &nodeInformer{
		client:       client,
		factory:      factory,
		listenerChan: make(chan cache.ResourceEventHandler),
		store:        cache.NewStore(nodeKeyFunc),
	}
}

var nodeKeyFunc = func(obj interface{}) (string, error) {
	node := obj.(*apicorev1.Node)
	return node.Name, nil
}

func (n *nodeInformer) List(selector labels.Selector) (ret []*apicorev1.Node, err error) {
	list, err := n.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}
	ret = make([]*apicorev1.Node, 0, 10)
	for i := 0; i < len(list.Items); i++ {
		node := list.Items[i]
		if selector.Matches(labels.Set(node.Labels)) {
			ret = append(ret, &node)
		}
	}
	return
}

func (n *nodeInformer) Get(name string) (*apicorev1.Node, error) {
	return n.client.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
}

func (n *nodeInformer) defaultInformer(client kubernetes.Interface, resync time.Duration) cache.SharedIndexInformer {
	informer, err := NewSharedIndexInformer(util.TopicNode, nodeKeyFunc)
	if err != nil {
		panic("无法创建Node的SharedIndexInformer")
	}

	return informer
}

func (n *nodeInformer) Informer() cache.SharedIndexInformer {
	return n.factory.InformerFor(&apicorev1.Node{}, n.defaultInformer)
}

func (n *nodeInformer) Lister() listerv1.NodeLister {
	return n
}
