package informers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/admissionregistration"
	"k8s.io/client-go/informers/apps"
	"k8s.io/client-go/informers/auditregistration"
	"k8s.io/client-go/informers/autoscaling"
	"k8s.io/client-go/informers/batch"
	"k8s.io/client-go/informers/certificates"
	"k8s.io/client-go/informers/coordination"
	"k8s.io/client-go/informers/core"
	"k8s.io/client-go/informers/discovery"
	"k8s.io/client-go/informers/events"
	"k8s.io/client-go/informers/extensions"
	"k8s.io/client-go/informers/flowcontrol"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/informers/networking"
	"k8s.io/client-go/informers/node"
	"k8s.io/client-go/informers/policy"
	"k8s.io/client-go/informers/rbac"
	"k8s.io/client-go/informers/scheduling"
	"k8s.io/client-go/informers/settings"
	"k8s.io/client-go/informers/storage"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"sync"
)

func NewSharedInformerFactory(client kubernetes.Interface) informers.SharedInformerFactory {
	return &sharedInformerFactory{
		client:           client,
		informers:        make(map[reflect.Type]cache.SharedIndexInformer),
		startedInformers: make(map[reflect.Type]bool),
	}
}

type sharedInformerFactory struct {
	started          bool
	informers        map[reflect.Type]cache.SharedIndexInformer
	startedInformers map[reflect.Type]bool
	client           kubernetes.Interface
	lock             sync.Mutex
}

func (f *sharedInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for resourceType, informer := range f.informers {
		if !f.startedInformers[resourceType] {
			go informer.Run(stopCh)
			f.startedInformers[resourceType] = true
		}
	}
}

func (f *sharedInformerFactory) InformerFor(obj runtime.Object, newFunc internalinterfaces.NewInformerFunc) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	typ := reflect.TypeOf(obj)
	var informer cache.SharedIndexInformer
	if informer, ok := f.informers[typ]; ok {
		return informer
	}

	informer = newFunc(f.client, 0)
	f.informers[typ] = informer
	return informer
}

func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (informers.GenericInformer, error) {
	panic("implement me")
}

func (f *sharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool {
	panic("implement me")
}

func (f *sharedInformerFactory) Admissionregistration() admissionregistration.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Apps() apps.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Auditregistration() auditregistration.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Autoscaling() autoscaling.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Batch() batch.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Certificates() certificates.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Coordination() coordination.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Core() core.Interface {
	return &coreInformer{
		client:  f.client,
		factory: f,
	}
}

func (f *sharedInformerFactory) Discovery() discovery.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Events() events.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Extensions() extensions.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Flowcontrol() flowcontrol.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Networking() networking.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Node() node.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Policy() policy.Interface {
	return &policyInformer{}
}

func (f *sharedInformerFactory) Rbac() rbac.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Scheduling() scheduling.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Settings() settings.Interface {
	panic("implement me")
}

func (f *sharedInformerFactory) Storage() storage.Interface {
	return &storageInformer{}
}
