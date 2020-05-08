package informers

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

func NewSharedIndexInformer(watcher watch.Interface, keyFunc cache.KeyFunc) (cache.SharedIndexInformer, error) {
	return &sharedIndexInformer{
		watcher:      watcher,
		listenerChan: make(chan cache.ResourceEventHandler),
		store:        cache.NewStore(keyFunc),
		isStop:       false,
	}, nil
}

type sharedIndexInformer struct {
	watcher       watch.Interface
	store         cache.Store
	listenerChan  chan cache.ResourceEventHandler
	channelClosed bool
	isStop        bool
}

func (s *sharedIndexInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	s.AddEventHandlerWithResyncPeriod(handler, 0)
}

func (s *sharedIndexInformer) AddEventHandlerWithResyncPeriod(handler cache.ResourceEventHandler, resyncPeriod time.Duration) {
	if s.isStop {
		if !s.channelClosed {
			close(s.listenerChan)
			s.channelClosed = true
		}
		// 停止后将不再接收任何
		return
	}

	s.listenerChan <- handler
}

func (s *sharedIndexInformer) GetStore() cache.Store {
	return s.store
}

func (s *sharedIndexInformer) GetController() cache.Controller {
	return s
}

func (s *sharedIndexInformer) Run(stopCh <-chan struct{}) {
	listeners := make([]cache.ResourceEventHandler, 0, 10)
	resultChan := s.watcher.ResultChan()
	for {
		logrus.Trace("Waiting to receive")
		select {
		case <-stopCh:
			logrus.Debug("Received stop signal, exiting")

			s.watcher.Stop()
			s.isStop = true
			return
		case ev := <-resultChan:
			switch ev.Type {
			case watch.Added:
				logrus.Debugf("Received added event, with object: %v", ev.Object)

				// 添加深拷贝对象到缓存
				err := s.store.Add(ev.Object.DeepCopyObject())
				if err != nil {
					logrus.Errorf("Error adding object %v to store: %v", ev.Object, err)
					continue
				}

				// 通知各个监听器
				for i := 0; i < len(listeners); i++ {
					listeners[i].OnAdd(ev.Object)
				}
			case watch.Modified:
				logrus.Debugf("Received modified event, with object: %v", ev.Object)

				key, err := nodeKeyFunc(ev.Object)
				if err != nil {
					logrus.Errorf("Error getting key for object %v : %v", ev.Object, err)
				}
				oldVal, exists, err := s.store.GetByKey(key)
				if !exists {
					logrus.Errorf("No object with key %s. This is abnormal because this is updating.", key)
					continue
				}
				if err != nil {
					logrus.Errorf("Error getting object with key %s", key)
				}

				err = s.store.Update(ev.Object.DeepCopyObject())
				if err != nil {
					logrus.Errorf("Error storing object %v", ev.Object)
				}

				for _, listener := range listeners {
					listener.OnUpdate(oldVal, ev.Object)
				}
			case watch.Deleted:
				logrus.Debugf("Received deleted event, with object: %v", ev.Object)

				err := s.store.Delete(ev.Object)
				if err != nil {
					logrus.Errorf("Error deleting object %v", ev.Object)
				}

				for _, listener := range listeners {
					listener.OnDelete(ev.Object)
				}
			case watch.Error:
				logrus.Errorf("Received Error Event: %v", ev.Object)
			}
		case listener := <-s.listenerChan:
			logrus.Debug("Added new listener")
			listeners = append(listeners, listener)
		}
	}
}

func (s *sharedIndexInformer) HasSynced() bool {
	return true
}

func (s *sharedIndexInformer) LastSyncResourceVersion() string {
	return ""
}

func (s *sharedIndexInformer) AddIndexers(indexers cache.Indexers) error {
	panic("implement me")
}

func (s *sharedIndexInformer) GetIndexer() cache.Indexer {
	panic("implement me")
}
