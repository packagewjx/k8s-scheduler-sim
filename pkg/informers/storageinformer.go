package informers

import (
	v13 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers/storage/v1"
	"k8s.io/client-go/informers/storage/v1alpha1"
	"k8s.io/client-go/informers/storage/v1beta1"
	v12 "k8s.io/client-go/listers/storage/v1"
	"k8s.io/client-go/tools/cache"
)

// storageInformer nil informer
type storageInformer struct {
}

func (s *storageInformer) CSIDrivers() v1.CSIDriverInformer {
	panic("implement me")
}

func (s *storageInformer) CSINodes() v1.CSINodeInformer {
	return &csiNodesInformer{}
}

func (s *storageInformer) StorageClasses() v1.StorageClassInformer {
	return &storageClassInformer{}
}

func (s *storageInformer) VolumeAttachments() v1.VolumeAttachmentInformer {
	panic("implement me")
}

func (s *storageInformer) V1() v1.Interface {
	return s
}

func (s *storageInformer) V1alpha1() v1alpha1.Interface {
	panic("implement me")
}

func (s *storageInformer) V1beta1() v1beta1.Interface {
	panic("implement me")
}

// storageClassInformer nil informer
type storageClassInformer struct {
}

func (s *storageClassInformer) List(selector labels.Selector) (ret []*v13.StorageClass, err error) {
	panic("implement me")
}

func (s *storageClassInformer) Get(name string) (*v13.StorageClass, error) {
	panic("implement me")
}

func (s *storageClassInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (s *storageClassInformer) Lister() v12.StorageClassLister {
	return s
}

type csiNodesInformer struct {
}

func (c *csiNodesInformer) List(selector labels.Selector) (ret []*v13.CSINode, err error) {
	panic("implement me")
}

func (c *csiNodesInformer) Get(name string) (*v13.CSINode, error) {
	panic("implement me")
}

func (c *csiNodesInformer) Informer() cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (c *csiNodesInformer) Lister() v12.CSINodeLister {
	return c
}
