package forkjoin

import (
	"fmt"
	"testing"
)

type SumTask struct {
	arr   []int
	start int
	end   int
}

func (task *SumTask) Fork() (tasks []Task) {
	mid := (task.start + task.end) / 2
	return []Task{
		&SumTask{
			arr:   task.arr,
			start: task.start,
			end:   mid,
		},
		&SumTask{
			arr:   task.arr,
			start: mid,
			end:   task.end,
		},
	}
}

func (task *SumTask) IsLeaf() bool {
	return task.end-task.start <= 16
}

func (task *SumTask) Leaf() interface{} {
	result := 0
	for i := task.start; i < task.end; i++ {
		result += task.arr[i]
	}
	return result
}

func (task *SumTask) Join(results []interface{}) interface{} {
	result := 0
	for i := 0; i < len(results); i++ {
		result += results[i].(int)
	}
	return result
}

func TestSimple(t *testing.T) {
	arr := make([]int, 1000000)
	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}
	task := &SumTask{
		arr:   arr,
		start: 0,
		end:   len(arr),
	}

	pool := NewSimpleForkJoinPool()
	result := pool.Execute(task).(int)
	fmt.Println(result)
}

func TestPerformance(t *testing.T) {

}
