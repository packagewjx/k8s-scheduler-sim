package metrics

type Aggregator interface {
	// Aggregate 将新的统计数据纳入到总统计
	Aggregate(tickMetrics *TickMetrics) *PeriodMetrics

	// Get 获取最新的统计数据
	Get() *PeriodMetrics
}

func NewAggregator() Aggregator {
	return &aggregator{
		count:    0,
		cpuSum:   0,
		cpu60:    newRingQueue(60),
		cpu300:   newRingQueue(300),
		cpu1500:  newRingQueue(1500),
		memSum:   0,
		mem60:    newRingQueue(60),
		mem300:   newRingQueue(300),
		mem1500:  newRingQueue(1500),
		loadSum:  0,
		load60:   newRingQueue(60),
		load300:  newRingQueue(300),
		load1500: newRingQueue(1500),
	}
}

type aggregator struct {
	count        int
	cpuSum       float64
	cpu60        *ringQueue
	cpu300       *ringQueue
	cpu1500      *ringQueue
	memSum       float64
	mem60        *ringQueue
	mem300       *ringQueue
	mem1500      *ringQueue
	loadSum      float64
	load60       *ringQueue
	load300      *ringQueue
	load1500     *ringQueue
	latestMetric *PeriodMetrics
}

func (a *aggregator) Get() *PeriodMetrics {
	return a.latestMetric
}

func (a *aggregator) Aggregate(tickMetrics *TickMetrics) *PeriodMetrics {
	a.count++
	a.cpuSum += tickMetrics.CpuUsage
	a.memSum += tickMetrics.MemUsage
	a.loadSum += tickMetrics.Load

	a.latestMetric = &PeriodMetrics{
		CpuUsageLastTick:           tickMetrics.CpuUsage,
		CpuUsageAverage:            a.cpuSum / float64(a.count),
		CpuUsageAverageIn60Ticks:   a.cpu60.add(tickMetrics.CpuUsage),
		CpuUsageAverageIn300Ticks:  a.cpu300.add(tickMetrics.CpuUsage),
		CpuUsageAverageIn1500Ticks: a.cpu1500.add(tickMetrics.CpuUsage),
		MemUsageLastTick:           tickMetrics.MemUsage,
		MemUsageAverage:            a.memSum / float64(a.count),
		MemUsageAverageIn60Ticks:   a.mem60.add(tickMetrics.MemUsage),
		MemUsageAverageIn300Ticks:  a.mem300.add(tickMetrics.MemUsage),
		MemUsageAverageIn1500Ticks: a.mem1500.add(tickMetrics.MemUsage),
		LoadLastTick:               tickMetrics.Load,
		LoadAverage:                a.loadSum / float64(a.count),
		LoadAverageIn60Ticks:       a.load60.add(tickMetrics.Load),
		LoadAverageIn300Ticks:      a.load300.add(tickMetrics.Load),
		LoadAverageIn1500Ticks:     a.load1500.add(tickMetrics.Load),
	}
	return a.latestMetric
}

func newRingQueue(capacity int) *ringQueue {
	return &ringQueue{
		arr:  make([]float64, capacity),
		size: 0,
		head: 0,
	}
}

type ringQueue struct {
	arr  []float64
	size int
	head int
	sum  float64
}

func (queue *ringQueue) add(num float64) (average float64) {
	if queue.size < len(queue.arr) {
		queue.arr[queue.size] = num
		queue.sum += num
		queue.size++
	} else /*queue.size == len(queue.arr)*/ {
		queue.sum -= queue.arr[queue.head]
		queue.sum += num
		queue.arr[queue.head] = num
		queue.head = (queue.head + 1) % queue.size
	}
	return queue.sum / float64(queue.size)
}
