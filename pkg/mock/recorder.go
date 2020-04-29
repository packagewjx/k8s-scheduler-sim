package mock

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
)

type eventRecorder struct {
}

func (recorder *eventRecorder) Eventf(regarding runtime.Object, related runtime.Object, eventtype, reason, action, note string, args ...interface{}) {
	s := fmt.Sprintf("regarding: %v. related: %v. eventtype: %s. reason: %s. action: %s. note: %s. args:", regarding, related, eventtype, reason, action, note)
	for i := 0; i < len(args); i++ {
		s += fmt.Sprintf(" %v,", args[i])
	}
	fmt.Println(s)
}

func SimRecorderFactory(name string) events.EventRecorder {
	return &eventRecorder{}
}
