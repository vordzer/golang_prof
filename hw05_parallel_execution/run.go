package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.

func Run(tasks []Task, n, m int) error {
	// Place your code here.
	bad := 0
	wg := sync.WaitGroup{}
	ch := make(chan int, n)
	cur := 0
	for _, t := range tasks {
		cur++
		wg.Add(1)
		go func(_t Task) {
			ok := _t()
			if ok != nil {
				ch <- 1
			} else {
				ch <- 0
			}
			wg.Done()
		}(t)
		if cur == n {
			v := <-ch
			bad += v
			cur--
		}
		if bad > m {
			break
		}
	}
	wg.Wait()
	close(ch)
	if bad > m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
