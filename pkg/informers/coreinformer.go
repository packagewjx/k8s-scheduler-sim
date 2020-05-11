package informers

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type coreInformer struct {
	client  kubernetes.Interface
	factory informers.SharedInformerFactory
}

func (c *coreInformer) ComponentStatuses() corev1.ComponentStatusInformer {
	panic("implement me")
}

func (c *coreInformer) ConfigMaps() corev1.ConfigMapInformer {
	panic("implement me")
}

func (c *coreInformer) Endpoints() corev1.EndpointsInformer {
	panic("implement me")
}

func (c *coreInformer) Events() corev1.EventInformer {
	panic("implement me")
}

func (c *coreInformer) LimitRanges() corev1.LimitRangeInformer {
	panic("implement me")
}

func (c *coreInformer) Namespaces() corev1.NamespaceInformer {
	panic("implement me")
}

func (c *coreInformer) Nodes() corev1.NodeInformer {
	return NewNodeInformer(c.client, c.factory)
}

func (c *coreInformer) PersistentVolumes() corev1.PersistentVolumeInformer {
	return &persistentVolumeInformer{}
}

func (c *coreInformer) PersistentVolumeClaims() corev1.PersistentVolumeClaimInformer {
	return &persistentVolumeClaimInformer{}
}

func (c *coreInformer) Pods() corev1.PodInformer {
	return newPodInformer(c.client, c.factory)
}

func (c *coreInformer) PodTemplates() corev1.PodTemplateInformer {
	panic("implement me")
}

func (c *coreInformer) ReplicationControllers() corev1.ReplicationControllerInformer {
	return &replicationControllerInformer{}
}

func (c *coreInformer) ResourceQuotas() corev1.ResourceQuotaInformer {
	panic("implement me")
}

func (c *coreInformer) Secrets() corev1.SecretInformer {
	panic("implement me")
}

func (c *coreInformer) Services() corev1.ServiceInformer {
	return &serviceInformer{}
}

func (c *coreInformer) ServiceAccounts() corev1.ServiceAccountInformer {
	panic("implement me")
}

func (c *coreInformer) V1() corev1.Interface {
	return c
}

// persistentVolumeClaimInformer nil informer
type persistentVolumeClaimInformer struct {
}

func (p *persistentVolumeClaimInformer) List(selector labels.Selector) (ret []*v1.PersistentVolumeClaim, err error) {
	panic("implement me")
}

func (p *persistentVolumeClaimInformer) PersistentVolumeClaims(namespace string) listerv1.PersistentVolumeClaimNamespaceLister {
	panic("implement me")
}

func (p *persistentVolumeClaimInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (p *persistentVolumeClaimInformer) Lister() listerv1.PersistentVolumeClaimLister {
	return p
}

// persistentVolumeInformer nil informer
type persistentVolumeInformer struct {
}

func (p *persistentVolumeInformer) List(selector labels.Selector) (ret []*v1.PersistentVolume, err error) {
	panic("implement me")
}

func (p *persistentVolumeInformer) Get(name string) (*v1.PersistentVolume, error) {
	panic("implement me")
}

func (p *persistentVolumeInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (p *persistentVolumeInformer) Lister() listerv1.PersistentVolumeLister {
	return p
}

// serviceInformer nil informer
type serviceInformer struct {
}

func (s *serviceInformer) Get(name string) (*v1.Service, error) {
	panic("implement me")
}

func (s *serviceInformer) List(selector labels.Selector) (ret []*v1.Service, err error) {
	return []*v1.Service{}, nil
}

func (s *serviceInformer) Services(namespace string) listerv1.ServiceNamespaceLister {
	return s
}

func (s *serviceInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (s *serviceInformer) Lister() listerv1.ServiceLister {
	return s
}

type replicationControllerInformer struct {
}

func (r *replicationControllerInformer) List(selector labels.Selector) (ret []*v1.ReplicationController, err error) {
	return []*v1.ReplicationController{}, nil
}

func (r *replicationControllerInformer) ReplicationControllers(namespace string) listerv1.ReplicationControllerNamespaceLister {
	panic("implement me")
}

func (r *replicationControllerInformer) GetPodControllers(pod *v1.Pod) ([]*v1.ReplicationController, error) {
	return []*v1.ReplicationController{}, nil
}

func (r *replicationControllerInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (r *replicationControllerInformer) Lister() listerv1.ReplicationControllerLister {
	return r
}
