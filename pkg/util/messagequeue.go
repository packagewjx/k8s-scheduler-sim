package util

import (
	"fmt"
	"k8s.io/apimachinery/pkg/watch"
	"sync"
)

const messageQueueSize = 16

type MessageQueue interface {
	// NewTopic 创建一个新的沟通话题，让订阅者和发布者进行沟通
	NewTopic(topic string) error

	Subscribe(topic string) (watch.Interface, error)

	Publish(topic string, event *watch.Event) error

	Shutdown()
}

var queue MessageQueue

func init() {
	queue = &messageQueueImpl{
		dispatchers: make(map[string]chan *operation),
		lock:        sync.Mutex{},
	}
}

func GetMessageQueue() MessageQueue {
	return queue
}

type messageQueueImpl struct {
	dispatchers map[string]chan *operation
	lock        sync.Mutex
}

var _ MessageQueue = &messageQueueImpl{}

func (queue *messageQueueImpl) Shutdown() {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	ev := &operation{
		what: opShutdown,
		arg:  nil,
	}
	for _, ch := range queue.dispatchers {
		ch <- ev
		close(ch)
	}
	queue.dispatchers = map[string]chan *operation{}
}

type operation struct {
	what opType
	arg  interface{}
}

type opType string

var (
	opAdd      opType = "add"
	opInform   opType = "inform"
	opDelete   opType = "delete"
	opShutdown opType = "shutdown"
)

func (queue *messageQueueImpl) NewTopic(topic string) error {
	if _, ok := queue.dispatchers[topic]; !ok {
		queue.lock.Lock()
		if _, ok := queue.dispatchers[topic]; !ok {
			ch := make(chan *operation)
			queue.dispatchers[topic] = ch
			go worker(ch)
		}
		queue.lock.Unlock()
	}
	return nil
}

func worker(ch chan *operation) {
	listeners := make([]chan watch.Event, 0, 10)
	for {
		op := <-ch

		switch op.what {
		case opAdd:
			listeners = append(listeners, op.arg.(chan watch.Event))
		case opDelete:
			p := len(listeners)
			toDelete := op.arg.([]chan watch.Event)
			toDeleteMap := make(map[chan watch.Event]bool)
			for i := 0; i < len(toDelete); i++ {
				toDeleteMap[toDelete[i]] = true
			}

			for i := 0; i < p; i++ {
				if toDeleteMap[listeners[i]] {
					temp := listeners[i]
					listeners[i] = listeners[p-1]
					listeners[p-1] = temp
					p--
					i--
				}
			}
			listeners = listeners[:p]
		case opInform:
			ev := op.arg.(*watch.Event)
			for i := 0; i < len(listeners); i++ {
				listeners[i] <- *ev
			}
		case opShutdown:
			for i := 0; i < len(listeners); i++ {
				close(listeners[i])
			}
			return
		}
	}

}

func (queue *messageQueueImpl) Subscribe(topic string) (watch.Interface, error) {
	if _, ok := queue.dispatchers[topic]; !ok {
		return nil, fmt.Errorf("没有话题%s", topic)
	}
	return &watcher{
		channels: make([]chan watch.Event, 0, 1),
		topic:    topic,
		queue:    queue,
	}, nil
}

func (queue *messageQueueImpl) Publish(topic string, event *watch.Event) error {
	if _, ok := queue.dispatchers[topic]; !ok {
		return fmt.Errorf("没有话题%s", topic)
	}

	op := &operation{
		what: opInform,
		arg:  event,
	}
	queue.dispatchers[topic] <- op
	return nil
}

func (queue *messageQueueImpl) newChannel(topic string) chan watch.Event {
	c := make(chan watch.Event, messageQueueSize)
	op := &operation{
		what: opAdd,
		arg:  c,
	}
	queue.dispatchers[topic] <- op
	return c
}

func (queue *messageQueueImpl) removeTopicChannels(topic string, chans []chan watch.Event) {
	op := &operation{
		what: opDelete,
		arg:  chans,
	}

	queue.dispatchers[topic] <- op
}

type watcher struct {
	channels []chan watch.Event
	topic    string
	queue    *messageQueueImpl
}

func (w *watcher) Stop() {
	w.queue.removeTopicChannels(w.topic, w.channels)
}

func (w *watcher) ResultChan() <-chan watch.Event {
	ch := w.queue.newChannel(w.topic)
	w.channels = append(w.channels, ch)
	return ch
}
