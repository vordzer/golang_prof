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
	wg := sync.WaitGroup{}
	ch := make(chan Task, len(tasks))
	mu := sync.Mutex{}
	wg.Add(n)
	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			res := 0
			for {
				mu.Lock()
				m -= res
				if m <= 0 {
					mu.Unlock()
					return
				}
				mu.Unlock()
				t, cok := <-ch
				if t == nil && !cok {
					return
				}
				ok := t()
				if ok != nil {
					res = 1
				} else {
					res = 0
				}
			}
		}()
	}
	wg.Wait()
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
