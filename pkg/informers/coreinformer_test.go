package informers

import (
	"context"
	"fmt"
	"github.com/packagewjx/k8s-scheduler-sim/pkg/util/fake"
	"github.com/sirupsen/logrus"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"testing"
	"time"
)

func TestNodeInformer(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	fakeClient := fake.NewFakeKubernetesInterface()
	factory := NewSharedInformerFactory(fakeClient)

	// 测试通知是否正常

	nodeInformer := factory.Core().V1().Nodes().Informer()
	stopCh := make(chan struct{})
	go nodeInformer.Run(stopCh)
	ch := make(chan *apicorev1.Node)
	nodeInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add event")
			ch <- obj.(*apicorev1.Node)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("update event")
			ch <- oldObj.(*apicorev1.Node)
			ch <- newObj.(*apicorev1.Node)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete event")
			ch <- obj.(*apicorev1.Node)
		},
	})

	nodeClient := fakeClient.CoreV1().Nodes()
	oldGenerateName := "test-1"
	testNode := &apicorev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:         "test",
			GenerateName: oldGenerateName,
		},
	}
	nodeClient.Create(context.TODO(), testNode, metav1.CreateOptions{})
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		// 异步通知有时间延迟，需要等待
		t.Error("创建通知失败")
	}

	newGenerateName := "test-100"
	testNode.GenerateName = newGenerateName
	nodeClient.Update(context.TODO(), testNode, metav1.UpdateOptions{})
	for i := 0; i < 2; i++ {
		select {
		case n := <-ch:
			if i == 0 {
				if n.GenerateName != oldGenerateName {
					t.Error("旧对象字段不对")
				}
			} else {
				fmt.Println("haha")
				if n.GenerateName != newGenerateName {
					t.Error("新对象字段不对")
				}
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("更新通知失败")
		}
	}

	nodeClient.UpdateStatus(context.TODO(), testNode, metav1.UpdateOptions{})
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("更新通知失败")
	}

	nodeClient.Delete(context.TODO(), "", metav1.DeleteOptions{})
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("删除通知失败")
	}

	stopCh <- struct{}{}
	fmt.Println("haha")
}
