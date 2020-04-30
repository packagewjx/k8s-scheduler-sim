package forkjoin

import "sync"

type simpleForkJoinPool struct {
}

func (pool *simpleForkJoinPool) Shutdown() {
	return
}

func (pool *simpleForkJoinPool) Submit(task Task, resultChan chan interface{}) {
	go func() {
		result := pool.Execute(task)
		resultChan <- result
	}()
}

func NewSimpleForkJoinPool() Pool {
	return &simpleForkJoinPool{}
}

func (pool *simpleForkJoinPool) Execute(task Task) interface{} {
	if task.IsLeaf() {
		return task.Leaf()
	} else {
		tasks := task.Fork()
		results := make([]interface{}, len(tasks))
		wg := sync.WaitGroup{}
		for i := 0; i < len(tasks); i++ {
			wg.Add(1)
			go func(idx int) {
				results[idx] = pool.Execute(tasks[idx])
				wg.Done()
			}(i)
		}
		wg.Wait()

		return task.Join(results)
	}
}
