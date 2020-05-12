package metrics

import (
	"math"
	"testing"
)

func TestAggregator(t *testing.T) {
	agg := NewAggregator()

	sum := float64(0)
	sum60 := float64(0)
	sum300 := float64(0)
	sum1500 := float64(0)
	met := &TickMetrics{
		CpuUsage: 0,
		MemUsage: 0,
		Load:     0,
	}

	for i := float64(0); i < 10000; i++ {
		met.Load = i * 0.1
		met.MemUsage = i * 0.1
		met.CpuUsage = i * 0.1

		sum += i * 0.1
		sum1500 += i * 0.1
		sum300 += i * 0.1
		sum60 += i * 0.1
		base1500 := float64(1500)
		base300 := float64(300)
		base60 := float64(60)
		if i >= 1500 {
			sum1500 -= 0.1 * (i - 1500)
		} else {
			base1500 = i + 1
		}
		if i >= 300 {
			sum300 -= 0.1 * (i - 300)
		} else {
			base300 = i + 1
		}
		if i >= 60 {
			sum60 -= 0.1 * (i - 60)
		} else {
			base60 = i + 1
		}

		data := agg.Aggregate(met)
		avg := sum / (1 + i)
		if !floatEquals(avg, data.CpuUsageAverage) || !floatEquals(avg, data.LoadAverage) || !floatEquals(avg, data.MemUsageAverage) {
			t.Errorf("总平均不对")
		}
		if !floatEquals(sum1500/base1500, data.CpuUsageAverageIn1500Ticks) ||
			!floatEquals(sum300/base300, data.CpuUsageAverageIn300Ticks) ||
			!floatEquals(sum60/base60, data.CpuUsageAverageIn60Ticks) {
			t.Errorf("长期平均不对")
		}
		if !floatEquals(sum1500/base1500, data.MemUsageAverageIn1500Ticks) ||
			!floatEquals(sum300/base300, data.MemUsageAverageIn300Ticks) ||
			!floatEquals(sum60/base60, data.MemUsageAverageIn60Ticks) {
			t.Errorf("长期平均不对")
		}
		if !floatEquals(sum1500/base1500, data.LoadAverageIn1500Ticks) ||
			!floatEquals(sum300/base300, data.LoadAverageIn300Ticks) ||
			!floatEquals(sum60/base60, data.LoadAverageIn60Ticks) {
			t.Errorf("长期平均不对")
		}
	}
}

func TestRingQueue(t *testing.T) {
	queue := newRingQueue(60)

	sum := float64(0)
	for i := float64(0); i < 60; i++ {
		avg := queue.add(i)
		sum += i

		if !floatEquals(avg, sum/(i+1)) {
			t.Errorf("平均数不对")
		}
	}

	for i := float64(0); i < 60; i++ {
		sum -= i
		sum += 60
		avg := queue.add(60)
		if !floatEquals(avg, sum/60) {
			t.Errorf("平均数不对")
		}
	}
}

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}
