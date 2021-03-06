package informers

import (
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sync"
	"time"
)

// Please Note there is a race condition where the creation of resultChan will happen after the first event publishing,
// causing lost of event. To prevent this, wait for some times after starting the factory or Run before using any
//kubernetes.Interface function in order to let the informer subscribe topic first.
func NewSharedIndexInformer(topic string, keyFunc cache.KeyFunc) (cache.SharedIndexInformer, error) {
	return &sharedIndexInformer{
		keyFunc: keyFunc,
		topic:   topic,
		store:   cache.NewStore(keyFunc),
		isStop:  false,
	}, nil
}

type sharedIndexInformer struct {
	keyFunc   cache.KeyFunc
	topic     string
	store     cache.Store
	listeners []cache.ResourceEventHandler
	isStop    bool
	lock      sync.Mutex
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
	if !s.isStop {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.listeners = append(s.listeners, handler)
		list := s.store.List()
		if len(list) > 0 {
			logrus.Debugf("Sending new listener all added resources")
			for _, item := range list {
				handler.OnAdd(item)
			}
		}
	}
}

func (s *sharedIndexInformer) GetStore() cache.Store {
	return s.store
}

func (s *sharedIndexInformer) GetController() cache.Controller {
	return s
}

func (s *sharedIndexInformer) Run(stopCh <-chan struct{}) {
	err := util.GetMessageQueue().SubscribeSynchronously(s.topic, func(ev *watch.Event) {
		if s.isStop {
			return
		}

		switch ev.Type {
		case watch.Added:
			logrus.Debugf("Received added event, with object: %v", ev.Object)

			// 添加深拷贝对象到缓存
			err := s.store.Add(ev.Object.DeepCopyObject())
			if err != nil {
				logrus.Errorf("Error adding object %v to store: %v", ev.Object, err)
				return
			}

			// 通知各个监听器
			s.lock.Lock()
			for _, listener := range s.listeners {
				logrus.Tracef("Notifying listener %p", &listener)
				listener.OnAdd(ev.Object)
			}
			s.lock.Unlock()
		case watch.Modified:
			logrus.Debugf("Received modified event, with object: %v", ev.Object)

			key, err := s.keyFunc(ev.Object)
			if err != nil {
				logrus.Errorf("Error getting key for object %v : %v", ev.Object, err)
			}
			oldVal, exists, err := s.store.GetByKey(key)
			if !exists {
				logrus.Errorf("No object with key %s. This is abnormal because this is updating.", key)
				return
			}
			if err != nil {
				logrus.Errorf("Error getting object with key %s", key)
			}

			err = s.store.Update(ev.Object.DeepCopyObject())
			if err != nil {
				logrus.Errorf("Error storing object %v", ev.Object)
			}

			s.lock.Lock()
			for _, listener := range s.listeners {
				logrus.Tracef("Notifying listener %p", &listener)
				listener.OnUpdate(oldVal, ev.Object)
			}
			s.lock.Unlock()
		case watch.Deleted:
			logrus.Debugf("Received deleted event, with object: %v", ev.Object)

			err := s.store.Delete(ev.Object)
			if err != nil {
				logrus.Errorf("Error deleting object %v", ev.Object)
				return
			}

			s.lock.Lock()
			for _, listener := range s.listeners {
				logrus.Tracef("Notifying listener %p", &listener)
				listener.OnDelete(ev.Object)
			}
			s.lock.Unlock()
		case watch.Error:
			logrus.Errorf("Received Error Event: %v", ev.Object)
		}
	})
	if err != nil {
		logrus.Errorf("error subscribing topic %s", s.topic)
	}

	go func() {
		<-stopCh
		s.isStop = true
	}()
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
