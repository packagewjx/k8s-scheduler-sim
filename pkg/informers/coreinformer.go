package informers

import (
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
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
	panic("implement me")
}

func (c *coreInformer) PersistentVolumeClaims() corev1.PersistentVolumeClaimInformer {
	panic("implement me")
}

func (c *coreInformer) Pods() corev1.PodInformer {
	return newPodInformer(c.client, c.factory)
}

func (c *coreInformer) PodTemplates() corev1.PodTemplateInformer {
	panic("implement me")
}

func (c *coreInformer) ReplicationControllers() corev1.ReplicationControllerInformer {
	panic("implement me")
}

func (c *coreInformer) ResourceQuotas() corev1.ResourceQuotaInformer {
	panic("implement me")
}

func (c *coreInformer) Secrets() corev1.SecretInformer {
	panic("implement me")
}

func (c *coreInformer) Services() corev1.ServiceInformer {
	panic("implement me")
}

func (c *coreInformer) ServiceAccounts() corev1.ServiceAccountInformer {
	panic("implement me")
}

func (c *coreInformer) V1() corev1.Interface {
	return c
}
