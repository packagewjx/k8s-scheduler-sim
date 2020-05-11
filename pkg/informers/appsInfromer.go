package informers

import (
	v13 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/informers/apps/v1beta1"
	"k8s.io/client-go/informers/apps/v1beta2"
	v12 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

type appsInformer struct {
}

func (a *appsInformer) ControllerRevisions() v1.ControllerRevisionInformer {
	panic("implement me")
}

func (a *appsInformer) DaemonSets() v1.DaemonSetInformer {
	panic("implement me")
}

func (a *appsInformer) Deployments() v1.DeploymentInformer {
	panic("implement me")
}

func (a *appsInformer) ReplicaSets() v1.ReplicaSetInformer {
	return &replicaSetInformer{}
}

func (a *appsInformer) StatefulSets() v1.StatefulSetInformer {
	return &statefulSetInformer{}
}

func (a *appsInformer) V1() v1.Interface {
	return a
}

func (a *appsInformer) V1beta1() v1beta1.Interface {
	panic("implement me")
}

func (a *appsInformer) V1beta2() v1beta2.Interface {
	panic("implement me")
}

type replicaSetInformer struct {
}

func (r *replicaSetInformer) Get(name string) (*v13.ReplicaSet, error) {
	panic("implement me")
}

func (r *replicaSetInformer) List(selector labels.Selector) (ret []*v13.ReplicaSet, err error) {
	panic("implement me")
}

func (r *replicaSetInformer) ReplicaSets(namespace string) v12.ReplicaSetNamespaceLister {
	return r
}

func (r *replicaSetInformer) GetPodReplicaSets(pod *apicorev1.Pod) ([]*v13.ReplicaSet, error) {
	return []*v13.ReplicaSet{}, nil
}

func (r *replicaSetInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (r *replicaSetInformer) Lister() v12.ReplicaSetLister {
	return r
}

type statefulSetInformer struct {
}

func (s *statefulSetInformer) List(selector labels.Selector) (ret []*v13.StatefulSet, err error) {
	panic("implement me")
}

func (s *statefulSetInformer) StatefulSets(namespace string) v12.StatefulSetNamespaceLister {
	panic("implement me")
}

func (s *statefulSetInformer) GetPodStatefulSets(pod *apicorev1.Pod) ([]*v13.StatefulSet, error) {
	return []*v13.StatefulSet{}, nil
}

func (s *statefulSetInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (s *statefulSetInformer) Lister() v12.StatefulSetLister {
	return s
}
