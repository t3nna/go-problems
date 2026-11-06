package main

import (
	"context"
	"fmt"
	"sync"
)

// ============================================================================
// TASK 1: The Resilient Worker Pool (with Graceful Shutdown & Early Exit)
// ============================================================================

type Job struct {
	ID      int
	Execute func(ctx context.Context) error
}

// RunWorkerPool limits concurrent processing to exactly 'workers'.
// Requirements:
// 1. If any job returns an error, immediately signal all other workers to stop picking up new jobs.
// 2. Wait for all currently executing workers to finish before returning.
// 3. Prevent any goroutine or channel leaks.
func RunWorkerPool(ctx context.Context, workers int, jobs <-chan Job) error {
	if workers <= 0 {
		return fmt.Errorf("workers should be greater then 0")
	}

	pCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var once sync.Once
	var pErr error

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case w, ok := <-jobs:
					if !ok {
						return
					}
					if pCtx.Err() != nil {
						return
					}
					err := w.Execute(pCtx)
					if err != nil {
						once.Do(func() {
							pErr = err
						})

						cancel()
						return
					}

				case <-pCtx.Done():
					return
				}

			}
		}()

	}

	wg.Wait()

	if pErr != nil {
		return pErr
	}

	return ctx.Err()

}
