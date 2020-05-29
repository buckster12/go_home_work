package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrNegativeNumber = errors.New("negative number of go-routines")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, goroutinesNumber int, errorsLimit int) error {
	wg := sync.WaitGroup{}
	var totalMutex sync.Mutex
	totalErrors := 0

	if goroutinesNumber <= 0 {
		return ErrNegativeNumber
	}

	// A channel to limit go-routines
	concurrentGoroutines := make(chan struct{}, goroutinesNumber)

	// Loop all the tasks
	for _, task := range tasks {
		// If number of errors is over limit - exit
		if errorsLimit > 0 {
			totalMutex.Lock()
			if totalErrors >= errorsLimit {
				totalMutex.Unlock()
				break
			}
			totalMutex.Unlock()
		}

		// Run go-routine
		wg.Add(1)
		concurrentGoroutines <- struct{}{}
		go func(task Task) {
			defer wg.Done()

			ok := task()

			// Save error
			if ok != nil && errorsLimit > 0 {
				totalMutex.Lock()
				totalErrors++
				totalMutex.Unlock()
			}
			// free place for next go-routine
			<-concurrentGoroutines
		}(task)
	}

	// wait while have go-routines
	wg.Wait()

	if totalErrors >= errorsLimit && errorsLimit > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
