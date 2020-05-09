package informers

import (
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

type fakeInformer struct {
}

var newFakeInformer internalinterfaces.NewInformerFunc = func(_ kubernetes.Interface, _ time.Duration) cache.SharedIndexInformer {
	return &fakeInformer{}
}

func (f *fakeInformer) AddEventHandler(handler cache.ResourceEventHandler) {

}

func (f *fakeInformer) AddEventHandlerWithResyncPeriod(handler cache.ResourceEventHandler, resyncPeriod time.Duration) {

}

func (f *fakeInformer) GetStore() cache.Store {
	return nil
}

func (f *fakeInformer) GetController() cache.Controller {
	return f
}

func (f *fakeInformer) Run(stopCh <-chan struct{}) {
	<-stopCh
}

func (f *fakeInformer) HasSynced() bool {
	return true
}

func (f *fakeInformer) LastSyncResourceVersion() string {
	return ""
}

func (f *fakeInformer) AddIndexers(indexers cache.Indexers) error {
	return nil
}

func (f *fakeInformer) GetIndexer() cache.Indexer {
	return nil
}
