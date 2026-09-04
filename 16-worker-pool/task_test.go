package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestRunWorkerPool_Success verifies that when all jobs complete perfectly,
// the function returns nil and all jobs are processed.
func TestRunWorkerPool_Success(t *testing.T) {
	var processedCount int64
	jobCount := 10
	workers := 3

	jobsChan := make(chan Job, jobCount)
	for i := 1; i <= jobCount; i++ {
		jobsChan <- Job{
			ID: i,
			Execute: func(ctx context.Context) error {
				atomic.AddInt64(&processedCount, 1)
				return nil
			},
		}
	}
	close(jobsChan)

	err := RunWorkerPool(context.Background(), workers, jobsChan)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if atomic.LoadInt64(&processedCount) != int64(jobCount) {
		t.Errorf("expected %d jobs to be processed, but got %d", jobCount, processedCount)
	}
}

// TestRunWorkerPool_EarlyExit verifies that if a job returns an error,
// that exact error is returned, and workers stop picking up new jobs.
func TestRunWorkerPool_EarlyExit(t *testing.T) {
	expectedErr := errors.New("database timeout simulation")
	workers := 2

	// We load plenty of jobs into the buffer
	jobCount := 50
	jobsChan := make(chan Job, jobCount)

	var processedCount int64

	// Job 1 will fail immediately
	jobsChan <- Job{
		ID: 1,
		Execute: func(ctx context.Context) error {
			atomic.AddInt64(&processedCount, 1)
			return expectedErr
		},
	}

	// The rest of the jobs are slow filler jobs
	for i := 2; i <= jobCount; i++ {
		jobsChan <- Job{
			ID: i,
			Execute: func(ctx context.Context) error {
				atomic.AddInt64(&processedCount, 1)
				// Small sleep to simulate work and allow cancellation to propagate
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		}
	}
	close(jobsChan)

	err := RunWorkerPool(context.Background(), workers, jobsChan)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error '%v', got: %v", expectedErr, err)
	}

	// Because of the early exit, the pool should have canceled before processing all 50 jobs
	finalCount := atomic.LoadInt64(&processedCount)
	if finalCount == int64(jobCount) {
		t.Errorf("expected worker pool to exit early, but all %d jobs were executed", jobCount)
	}
}

// TestRunWorkerPool_ParentContextCanceled verifies that if the parent context
// is canceled externally, the pool gracefully halts execution.
func TestRunWorkerPool_ParentContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	workers := 2
	jobsChan := make(chan Job, 10)

	// Add a few infinite/blocking jobs
	for i := 1; i <= 5; i++ {
		jobsChan <- Job{
			ID: i,
			Execute: func(ctx context.Context) error {
				// Block until the worker pool's internal context tells us to stop
				<-ctx.Done()
				return ctx.Err()
			},
		}
	}
	close(jobsChan)

	// Cancel the parent context shortly after starting
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := RunWorkerPool(ctx, workers, jobsChan)

	// Since the parent context was canceled and no job explicitly failed with a custom error,
	// the returned error should either be nil (if treated as a clean shutdown) or the context error.
	// Based on your specific logic, it returns firstErr which remains nil.
	if err != nil {
		t.Errorf("expected clean exit or handled context, got: %v", err)
	}
}
