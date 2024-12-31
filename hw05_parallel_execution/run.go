package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorCount struct {
	mu          *sync.RWMutex
	errorNumber int
	err         error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.

func Run(tasks []Task, n, m int) error {
	taskChan := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}

	addTasks(taskChan, tasks)

	errCount := errorCount{
		err:         nil,
		mu:          &sync.RWMutex{},
		errorNumber: 0,
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go doTasks(taskChan, m, &wg, &errCount)
	}

	wg.Wait()

	return errCount.err
}

func addTasks(taskChan chan Task, tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		taskChan <- tasks[i]
	}
	close(taskChan)
}

func doTasks(taskChan chan Task, maxErrorNumber int, wg *sync.WaitGroup, errCount *errorCount) {
	defer wg.Done()
	currentTask, ok := <-taskChan
	for ok {
		errCount.mu.RLock()
		if errCount.errorNumber > maxErrorNumber && maxErrorNumber > 0 {
			errCount.mu.RUnlock()
			return
		}
		errCount.mu.RUnlock()
		err := currentTask()
		currentTask, ok = <-taskChan
		if maxErrorNumber > 0 && err != nil {
			errCount.mu.Lock()
			errCount.errorNumber++
			if errCount.errorNumber >= maxErrorNumber {
				errCount.err = ErrErrorsLimitExceeded
				errCount.mu.Unlock()
				return
			}
			errCount.mu.Unlock()
		}
	}
}
