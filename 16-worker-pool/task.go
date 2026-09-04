package main

import (
	"context"
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

	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	var errOnce sync.Once
	var firstErr error

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case val, ok := <-jobs:
					if !ok {
						return
					}
					if err := val.Execute(ctxC); err != nil {
						errOnce.Do(func() {
							firstErr = err
						})

						cancel()

					}
				case <-ctxC.Done():
					return

				}

			}
		}()
	}

	wg.Wait()

	return firstErr

}
