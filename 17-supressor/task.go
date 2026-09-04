package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// TASK 3: The "Thundering Herd" Suppressor (Custom Singleflight)
// ============================================================================

type call struct {
	val any
	err error
	wg  sync.WaitGroup
}
type Suppressor struct {
	// TODO: Add internal fields (hint: you will need a mutex and a map tracking active calls/channels)
	mu sync.Mutex
	m  map[string]*call
}

func NewSuppressor() *Suppressor {
	return &Suppressor{
		m: make(map[string]*call), // Initialize the map to prevent nil panic
	}
}

// Do executes 'fn' exactly once for concurrent requests using the same 'key'.
// Requirements:
// 1. Duplicate concurrent callers must block, wait for the first execution to finish, and get the same result/error.
// 2. Once the execution finishes, subsequent new calls for the key must execute 'fn' again.
// 3. Do not use the 'golang.org/x/sync/singleflight' package.

func (s *Suppressor) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	s.mu.Lock()

	// CASE A: The key is already in flight (You are a passenger)
	if c, exists := s.m[key]; exists {
		// CRITICAL: Unlock immediately so other keys aren't blocked while we wait
		s.mu.Unlock()

		// Block until the leader goroutine calls wg.Done()
		c.wg.Wait()

		// Return the results fetched by the leader
		return c.val, c.err
	}

	// CASE B: You are the first one here (You are the leader)
	c := &call{}
	c.wg.Add(1)
	s.m[key] = c

	// CRITICAL: Unlock before running the slow function so other goroutines
	// searching for this key (or other keys) can enter the map.
	s.mu.Unlock()

	// Execute the actual heavy/slow function outside the lock
	c.val, c.err = fn()

	// Execution finished. We must clean up the map.
	s.mu.Lock()
	// Delete the key so that FUTURE requests (after this flight completes)
	// will trigger a fresh execution of fn() instead of getting stale data.
	delete(s.m, key)
	s.mu.Unlock()

	// Broadcast to all waiting passengers that the data is ready
	c.wg.Done()

	return c.val, c.err
}

func main() {
	suppressor := NewSuppressor()

	var dbQueryCounter int32
	var wg sync.WaitGroup

	// This simulates our heavy database query or expensive API call
	slowDatabaseQuery := func() (interface{}, error) {
		// Increment the counter every time this function is ACTUALLY executed
		atomic.AddInt32(&dbQueryCounter, 1)

		// Simulate a 500ms network delay
		time.Sleep(500 * time.Millisecond)
		return "Product_Data_XYZ", nil
	}

	totalConcurrentRequests := 10
	fmt.Printf("--- Phase 1: Launching %d concurrent requests simultaneously ---\n", totalConcurrentRequests)

	wg.Add(totalConcurrentRequests)
	for i := 1; i <= totalConcurrentRequests; i++ {
		go func(workerID int) {
			defer wg.Done()

			// All workers request the exact same key ("product_42") at the same time
			result, err := suppressor.Do("product_42", slowDatabaseQuery)
			if err != nil {
				fmt.Printf("Worker %d failed: %v\n", workerID, err)
				return
			}
			fmt.Printf("Worker %d completed. Got result: %v\n", workerID, result)
		}(i)
	}

	// Wait for the first wave of thundering herd requests to complete
	wg.Wait()

	fmt.Println("\n--- Phase 1 Results ---")
	fmt.Printf("Total Goroutines served: %d\n", totalConcurrentRequests)
	fmt.Printf("Actual Database executions: %d (Should be 1)\n", atomic.LoadInt32(&dbQueryCounter))

	// -------------------------------------------------------------------------
	// Phase 2: Verify that future requests execute the function again (Requirement 2)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- Phase 2: Making a subsequent request after flight cleared ---")

	result, _ := suppressor.Do("product_42", slowDatabaseQuery)

	fmt.Println("--- Phase 2 Results ---")
	fmt.Printf("Subsequent request got result: %v\n", result)
	fmt.Printf("Total Database executions now: %d (Should be 2)\n", atomic.LoadInt32(&dbQueryCounter))
}
