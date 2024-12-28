package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorCount struct {
	mu          *sync.Mutex
	errorNumber int
	err         error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.

func Run(tasks []Task, n, m int) error {
	taskChan := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go addTasks(&taskChan, tasks, &wg)

	errCount := errorCount{
		err:         nil,
		mu:          &sync.Mutex{},
		errorNumber: 0,
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go doTasks(&taskChan, m, &wg, &errCount)
	}

	wg.Wait()

	return errCount.err
}

func addTasks(taskChan *chan Task, tasks []Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(tasks); i++ {
		*taskChan <- tasks[i]
	}
	close(*taskChan)
}

func doTasks(taskChan *chan Task, maxErrorNumber int, wg *sync.WaitGroup, errCount *errorCount) {
	defer wg.Done()
	currentTask, ok := <-*taskChan
	for ok {
		if errCount.errorNumber < maxErrorNumber || maxErrorNumber <= 0 {
			err := currentTask()
			currentTask, ok = <-*taskChan
			if maxErrorNumber > 0 && err != nil {
				errCount.mu.Lock()
				errCount.errorNumber++
				if errCount.errorNumber >= maxErrorNumber {
					errCount.err = ErrErrorsLimitExceeded
					ok = false
				}
				errCount.mu.Unlock()
			}
		}
	}
}
