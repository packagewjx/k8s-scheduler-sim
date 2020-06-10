package controllers

import (
	rand2 "math/rand"
	"testing"
	"time"
)

func TestFindPercentile(t *testing.T) {
	arr := make([]uint8, 0, 100)
	for i := 0; i < 100; i++ {
		arr = append(arr, uint8(i))
	}

	// 随机打乱arr
	rand2.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {
		p := rand2.Intn(100 - i)
		temp := arr[p]
		arr[p] = arr[100-i-1]
		arr[100-i-1] = temp
	}

	for i := uint8(0); i < 100; i++ {
		if i != findPercentile(arr, float64(i)) {
			t.Errorf("wrong")
		}
	}
}

func TestServiceMetrics(t *testing.T) {
	m := &serviceMetrics{}
	for i := 0; i < 100; i++ {
		m.add(uint8(i))
	}

	avg, p50, p90, p99 := m.getMetrics()
	if p50 != 50 || p90 != 90 || p99 != 99 {
		t.Errorf("wrong")
	}
	t.Log(avg)
}
