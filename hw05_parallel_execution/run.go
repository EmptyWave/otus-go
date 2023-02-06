package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
	ErrInvalidNumberWorkers = errors.New("invalid number workers")
)

type Task func() error

type ErrorCount struct {
	value int
	mu    sync.RWMutex
}

func (err *ErrorCount) Inc() {
	err.mu.Lock()
	defer err.mu.Unlock()

	err.value++
}

func (err *ErrorCount) Get() int {
	err.mu.RLock()
	defer err.mu.RUnlock()

	return err.value
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m < 0 {
		return ErrInvalidNumberWorkers
	}

	errCount := ErrorCount{}
	workers := make(chan struct{}, n)
	wg := sync.WaitGroup{}
	defer close(workers)
	defer wg.Wait()

	for _, task := range tasks {
		if errCount.Get() >= m {
			return ErrErrorsLimitExceeded
		}

		workers <- struct{}{}
		wg.Add(1)
		go RunTusk(task, workers, &errCount, &wg)
	}

	return nil
}

func RunTusk(task Task, workers chan struct{}, errCount *ErrorCount, wg *sync.WaitGroup) {
	defer wg.Done()

	if task() != nil {
		errCount.Inc()
	}

	<-workers
}
