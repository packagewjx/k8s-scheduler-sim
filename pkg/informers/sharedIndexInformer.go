package informers

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sync"
	"time"
)

// Please Note there is a race condition where the creation of resultChan will happen after the first event publishing,
// causing lost of event. To prevent this, wait for some times after starting the factory or Run before using any
//kubernetes.Interface function in order to let the informer subscribe topic first.
func NewSharedIndexInformer(watcher watch.Interface, keyFunc cache.KeyFunc) (cache.SharedIndexInformer, error) {
	return &sharedIndexInformer{
		keyFunc:      keyFunc,
		watcher:      watcher,
		listenerChan: make(chan newHandlerEvent),
		store:        cache.NewStore(keyFunc),
		isStop:       false,
	}, nil
}

type sharedIndexInformer struct {
	keyFunc       cache.KeyFunc
	watcher       watch.Interface
	store         cache.Store
	listenerChan  chan newHandlerEvent
	listeners     []cache.ResourceEventHandler
	channelClosed bool
	isStarted     bool
	isStop        bool
	lock          sync.Mutex
}

func (s *sharedIndexInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	s.AddEventHandlerWithResyncPeriod(handler, 0)
}

type newHandlerEvent struct {
	handler cache.ResourceEventHandler
	// addedCh 控制添加监听器的同步性
	addedCh chan bool
}

func (s *sharedIndexInformer) AddEventHandlerWithResyncPeriod(handler cache.ResourceEventHandler, resyncPeriod time.Duration) {
	if s.isStop {
		if !s.channelClosed {
			close(s.listenerChan)
			s.channelClosed = true
		}
		// 停止后将不再接收任何新的监听器
		return
	}
	if s.isStarted {
		nh := newHandlerEvent{
			handler: handler,
			addedCh: make(chan bool),
		}
		s.listenerChan <- nh
	} else {
		s.lock.Lock()
		s.listeners = append(s.listeners, handler)
		logrus.Debugf("Added new listener %p", &handler)
		s.lock.Unlock()
	}

}

func (s *sharedIndexInformer) GetStore() cache.Store {
	return s.store
}

func (s *sharedIndexInformer) GetController() cache.Controller {
	return s
}

func (s *sharedIndexInformer) Run(stopCh <-chan struct{}) {
	logrus.Debug("SharedIndexInformer starting")

	s.isStarted = true
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
				for _, listener := range s.listeners {
					logrus.Tracef("Notifying listener %p", &listener)
					listener.OnAdd(ev.Object)
				}
			case watch.Modified:
				logrus.Debugf("Received modified event, with object: %v", ev.Object)

				key, err := s.keyFunc(ev.Object)
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

				for _, listener := range s.listeners {
					logrus.Tracef("Notifying listener %p", &listener)
					listener.OnUpdate(oldVal, ev.Object)
				}
			case watch.Deleted:
				logrus.Debugf("Received deleted event, with object: %v", ev.Object)

				err := s.store.Delete(ev.Object)
				if err != nil {
					logrus.Errorf("Error deleting object %v", ev.Object)
				}

				for _, listener := range s.listeners {
					logrus.Tracef("Notifying listener %p", &listener)
					listener.OnDelete(ev.Object)
				}
			case watch.Error:
				logrus.Errorf("Received Error Event: %v", ev.Object)
			}
		case listenerEvent := <-s.listenerChan:
			listener := listenerEvent.handler
			logrus.Debugf("Added new listener %p", &listener)
			s.listeners = append(s.listeners, listener)
			// notify that listener all added resource
			list := s.store.List()
			if len(list) > 0 {
				logrus.Debugf("Sending new listener all added resources")
				for _, item := range list {
					listener.OnAdd(item)
				}
			}
			listenerEvent.addedCh <- true
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
