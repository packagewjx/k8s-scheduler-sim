package informers

import (
	"fmt"
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
		select {
		case <-stopCh:
			s.watcher.Stop()
			s.isStop = true
			break
		case ev := <-resultChan:
			switch ev.Type {
			case watch.Added:
				// 添加节点到缓存
				err := s.store.Add(ev.Object)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// 通知各个监听器
				for i := 0; i < len(listeners); i++ {
					listeners[i].OnAdd(ev.Object)
				}
			case watch.Modified:
				key, _ := nodeKeyFunc(ev.Object)
				oldVal, exists, err := s.store.GetByKey(key)
				if !exists {
					fmt.Println("居然不存在")
					continue
				}
				if err != nil {
					fmt.Println(err)
				}

				err = s.store.Update(ev.Object)
				if err != nil {
					fmt.Println(err)
				}

				for _, listener := range listeners {
					listener.OnUpdate(oldVal, ev.Object)
				}
			case watch.Deleted:
				err := s.store.Delete(ev.Object)
				if err != nil {
					fmt.Println(err)
				}

				for _, listener := range listeners {
					listener.OnDelete(ev.Object)
				}
			case watch.Error:
				fmt.Printf("error:%v\n", ev.Object)
			}
		case listener := <-s.listenerChan:
			listeners = append(listeners, listener)
		}
		break
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
