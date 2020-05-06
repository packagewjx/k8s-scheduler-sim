package util

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/kubernetes/pkg/apis/core"
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	messageQueue := GetMessageQueue()
	recvTimes := 10

	topics := []string{"Animation", "Comics", "Game"}
	wg := sync.WaitGroup{}

	for i := 0; i < len(topics); i++ {
		_ = messageQueue.NewTopic(topics[i])
		wt, _ := messageQueue.Subscribe(topics[i])
		for j := 0; j < 10; j++ {
			resultChan := wt.ResultChan()
			wg.Add(1)
			go func(ch <-chan watch.Event) {
				defer wg.Done()
				for k := 0; k < recvTimes; k++ {
					<-ch
				}
			}(resultChan)
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < recvTimes; j++ {
			topic := topics[i]
			_ = messageQueue.Publish(topic, &watch.Event{
				Type:   watch.Added,
				Object: &core.Pod{},
			})
		}
	}

	wg.Wait()
}

func TestShutdown(t *testing.T) {

}
