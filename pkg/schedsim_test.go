package pkg

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestScheduleOne(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	sim := NewSchedulerSimulator()
	sim.Run()
}
