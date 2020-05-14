package util

import (
	"testing"
)

func TestRingQueue(t *testing.T) {
	queue := NewRingQueue(10)
	for i := 0; i < 100; i++ {
		queue.Offer(i)
	}

	cnt := 0
	queue.Do(func(val interface{}) {
		if cnt != val {
			t.Error("no correct")
		}
		cnt++
	})

	for i := 0; i < 100; i++ {
		if queue.Poll() != i {
			t.Error("no correct")
		}
	}

	queue = NewRingQueue(10)

	numAdd := 0
	numPoll := 0
	for ; numAdd < 10; numAdd++ {
		queue.Offer(numAdd)
	}

	for ; numPoll < 5; numPoll++ {
		if queue.Poll() != numPoll {
			t.Error("no correct")
		}
	}
	cnt = 5
	queue.Do(func(val interface{}) {
		if cnt != val {
			t.Error("no correct")
		}
		cnt++
	})

	for ; numAdd < 20; numAdd++ {
		queue.Offer(numAdd)
	}

	for ; numPoll < 15; numPoll++ {
		if queue.Poll() != numPoll {
			t.Error("no correct")
		}
	}
	cnt = 15
	queue.Do(func(val interface{}) {
		if cnt != val {
			t.Error("no correct")
		}
		cnt++
	})

	for ; numAdd < 30; numAdd++ {
		queue.Offer(numAdd)
	}
	for ; numPoll < 30; numPoll++ {
		if queue.Poll() != numPoll {
			t.Error("no correct")
		}
	}

	if queue.Poll() != nil {
		t.Error("no correct")
	}

}
