package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const testTimeout = 5 * time.Second

func waitSignal(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for signal")
	}
}

func waitResult(t *testing.T, ch <-chan error) error {
	t.Helper()
	select {
	case err := <-ch:
		return err
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for worker pool")
		return nil
	}
}

func TestRunWorkerPoolSuccess(t *testing.T) {
	const count = 10
	var processed atomic.Int32
	jobs := make(chan Job, count)
	for i := 0; i < count; i++ {
		jobs <- Job{ID: i, Execute: func(context.Context) error {
			processed.Add(1)
			return nil
		}}
	}
	close(jobs)
	if err := RunWorkerPool(context.Background(), 3, jobs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := processed.Load(); got != count {
		t.Errorf("processed %d jobs, want %d", got, count)
	}
}

func TestRunWorkerPoolConcurrencyLimit(t *testing.T) {
	const workers, count = 3, 7
	started := make(chan struct{}, count)
	release := make(chan struct{})
	var active, peak atomic.Int32
	jobs := make(chan Job, count)
	for i := 0; i < count; i++ {
		jobs <- Job{Execute: func(context.Context) error {
			n := active.Add(1)
			for old := peak.Load(); n > old && !peak.CompareAndSwap(old, n); old = peak.Load() {
			}
			started <- struct{}{}
			<-release
			active.Add(-1)
			return nil
		}}
	}
	close(jobs)
	result := make(chan error, 1)
	go func() { result <- RunWorkerPool(context.Background(), workers, jobs) }()
	for i := 0; i < workers; i++ {
		waitSignal(t, started)
	}
	close(release)
	if err := waitResult(t, result); err != nil {
		t.Fatal(err)
	}
	if got := peak.Load(); got != workers {
		t.Errorf("peak concurrency = %d, want %d", got, workers)
	}
}

func TestRunWorkerPoolJobErrorWaitsForRunningJob(t *testing.T) {
	expected := errors.New("job failed")
	started := make(chan struct{})
	canceled := make(chan struct{})
	release := make(chan struct{})
	var later atomic.Int32
	jobs := make(chan Job, 3)
	jobs <- Job{Execute: func(ctx context.Context) error {
		close(started)
		<-ctx.Done()
		close(canceled)
		<-release
		return nil
	}}
	jobs <- Job{Execute: func(context.Context) error {
		<-started
		return expected
	}}
	jobs <- Job{Execute: func(context.Context) error { later.Add(1); return nil }}
	close(jobs)
	result := make(chan error, 1)
	go func() { result <- RunWorkerPool(context.Background(), 2, jobs) }()
	waitSignal(t, canceled)
	select {
	case err := <-result:
		t.Fatalf("returned before running job finished: %v", err)
	default:
	}
	close(release)
	if err := waitResult(t, result); !errors.Is(err, expected) {
		t.Errorf("error = %v, want %v", err, expected)
	}
	if got := later.Load(); got != 0 {
		t.Errorf("started %d jobs after failure", got)
	}
}

func TestRunWorkerPoolParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	jobs := make(chan Job, 1)
	jobs <- Job{Execute: func(ctx context.Context) error {
		close(started)
		<-ctx.Done()
		return nil
	}}
	result := make(chan error, 1)
	go func() { result <- RunWorkerPool(ctx, 1, jobs) }()
	waitSignal(t, started)
	cancel()
	if err := waitResult(t, result); !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
}

func TestRunWorkerPoolCanceledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var executed atomic.Bool
	jobs := make(chan Job, 1)
	jobs <- Job{Execute: func(context.Context) error { executed.Store(true); return nil }}
	close(jobs)
	if err := RunWorkerPool(ctx, 1, jobs); !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
	if executed.Load() {
		t.Error("executed job after cancellation")
	}
}

func TestRunWorkerPoolEmptyJobs(t *testing.T) {
	jobs := make(chan Job)
	close(jobs)
	if err := RunWorkerPool(context.Background(), 2, jobs); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunWorkerPoolInvalidWorkers(t *testing.T) {
	for _, workers := range []int{0, -1} {
		if err := RunWorkerPool(context.Background(), workers, nil); err == nil {
			t.Errorf("workers = %d: expected an error", workers)
		}
	}
}

func TestRunWorkerPoolConcurrentErrors(t *testing.T) {
	const workers = 4
	start := make(chan struct{})
	var ready sync.WaitGroup
	ready.Add(workers)
	jobs := make(chan Job, workers)
	errs := make([]error, workers)
	for i := range errs {
		errs[i] = errors.New("job error")
		err := errs[i]
		jobs <- Job{Execute: func(context.Context) error {
			ready.Done()
			<-start
			return err
		}}
	}
	close(jobs)
	result := make(chan error, 1)
	go func() { result <- RunWorkerPool(context.Background(), workers, jobs) }()
	readyCh := make(chan struct{})
	go func() { ready.Wait(); close(readyCh) }()
	waitSignal(t, readyCh)
	close(start)
	got := waitResult(t, result)
	for _, err := range errs {
		if errors.Is(got, err) {
			return
		}
	}
	t.Errorf("error = %v, want one of the job errors", got)
}
