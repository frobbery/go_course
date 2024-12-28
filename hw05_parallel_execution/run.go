package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type taskIndexLocked struct {
	mu *sync.Mutex

	taskIndex int

	errorNumber int

	err error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.

func Run(tasks []Task, n, m int) error {
	taskIndex := taskIndexLocked{
		err: nil,

		mu: &sync.Mutex{},

		taskIndex: 0,

		errorNumber: 0,
	}

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go doTasks(tasks, m, &wg, &taskIndex)
	}

	wg.Wait()

	return taskIndex.err
}

func doTasks(tasks []Task, maxErrorNumber int, wg *sync.WaitGroup, taskIndex *taskIndexLocked) {
	defer wg.Done()

	currTaskIndex := getNextTaskIndex(taskIndex)

	for currTaskIndex < len(tasks) {
		task := tasks[currTaskIndex]

		err := task()

		if maxErrorNumber > 0 && err != nil {
			taskIndex.mu.Lock()

			taskIndex.errorNumber++

			if taskIndex.errorNumber >= maxErrorNumber {
				taskIndex.taskIndex = len(tasks)

				taskIndex.err = ErrErrorsLimitExceeded
			}

			taskIndex.mu.Unlock()
		}

		currTaskIndex = getNextTaskIndex(taskIndex)
	}
}

func getNextTaskIndex(taskIndex *taskIndexLocked) int {
	taskIndex.mu.Lock()

	index := taskIndex.taskIndex

	taskIndex.taskIndex++

	taskIndex.mu.Unlock()

	return index
}
