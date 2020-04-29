package metrics

type TickMetrics struct {
	// CpuUsage CPU使用百分比
	CpuUsage float64
	// MemUsage 内存使用百分比
	MemUsage float64
	Load     float64
}

type PeriodMetrics struct {
	CpuUsageAverage            float64
	CpuUsageAverageIn60Ticks   float64
	CpuUsageAverageIn300Ticks  float64
	CpuUsageAverageIn1500Ticks float64
	MemUsageAverage            float64
	MemUsageAverageIn60Ticks   float64
	MemUsageAverageIn300Ticks  float64
	MemUsageAverageIn1500Ticks float64
	LoadAverage                float64
	LoadAverageIn60Ticks       float64
	LoadAverageIn300Ticks      float64
	LoadAverageIn1500Ticks     float64
}

type ServiceCallMetrics struct {
	ElapsedTime float64
}

type ServiceMetrics struct {
	CallTimeAverage float64
	// 第99百分位的调用时间
	CallTime99th float64
	// 第99.99百分位的调动时间
	CallTime9999th float64
}
