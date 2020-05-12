package util

import (
	"github.com/sirupsen/logrus"
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
	queue = NewSynchronizedMessageQueue()
}

func GetMessageQueue() MessageQueue {
	return queue
}

var _ MessageQueue = &synchronizedMessageQueue{}

type watcher struct {
	channels            []chan watch.Event
	topic               string
	removeTopicChannels func(topic string, chans []chan watch.Event)
	newChannel          func(topic string) chan watch.Event
}

func (w *watcher) Stop() {
	w.removeTopicChannels(w.topic, w.channels)
}

func (w *watcher) ResultChan() <-chan watch.Event {
	ch := w.newChannel(w.topic)
	w.channels = append(w.channels, ch)
	return ch
}

// NewSynchronizedMessageQueue 构造一个全同步的消息队列。Subscribe方法将会在成功添加监听器后返回。Publish方法将会在成功
// 发送消息到所有监听器后返回。
func NewSynchronizedMessageQueue() MessageQueue {
	return &synchronizedMessageQueue{
		listeners: make(map[string][]chan watch.Event),
		lock:      sync.RWMutex{},
	}
}

type synchronizedMessageQueue struct {
	listeners map[string][]chan watch.Event
	lock      sync.RWMutex
}

func (queue *synchronizedMessageQueue) NewTopic(topic string) error {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	if _, ok := queue.listeners[topic]; !ok {
		queue.listeners[topic] = make([]chan watch.Event, 0, 10)
	}
	return nil
}

func (queue *synchronizedMessageQueue) Subscribe(topic string) (watch.Interface, error) {
	return &watcher{
		channels:            make([]chan watch.Event, 0, 10),
		topic:               topic,
		removeTopicChannels: queue.removeTopicChannels,
		newChannel:          queue.newChannel,
	}, nil
}

func (queue *synchronizedMessageQueue) Publish(topic string, event *watch.Event) error {
	queue.lock.RLock()
	defer queue.lock.RUnlock()

	logrus.Debugf("Publishing event %v", event)

	for _, listener := range queue.listeners[topic] {
		logrus.Tracef("Notifying subscriber %p", &listener)
		listener <- *event
	}

	logrus.Debugf("Finished notifying all subscribers.")
	return nil
}

func (queue *synchronizedMessageQueue) Shutdown() {
	return
}

func (queue *synchronizedMessageQueue) newChannel(topic string) chan watch.Event {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	ch := make(chan watch.Event)
	queue.listeners[topic] = append(queue.listeners[topic], ch)
	logrus.Debugf("Creating new subscriber %p", &ch)
	return ch
}

func (queue *synchronizedMessageQueue) removeTopicChannels(topic string, chans []chan watch.Event) {
	queue.lock.Lock()
	queue.lock.Unlock()

	listeners := queue.listeners[topic]
	p := len(listeners)
	toDeleteMap := make(map[chan watch.Event]bool)
	for i := 0; i < len(chans); i++ {
		toDeleteMap[chans[i]] = true
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
	queue.listeners[topic] = listeners[:p]
}
