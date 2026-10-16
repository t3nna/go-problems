package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func doWork(ctx context.Context, cancel context.CancelFunc, t int) {
	defer cancel()
	select {
	case <-time.After(time.Duration(t) * time.Second):
		fmt.Println("Work Done")
		cancel()

	case <-ctx.Done():
		fmt.Println("context canceled", ctx.Err())
	}
}

func doWorkCorrect(ctx context.Context, id, duration int, done chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-time.After(time.Duration(duration) * time.Second):
		fmt.Println("Worker is finished with id: ", id)
		select {
		case done <- id:
		default:
		}
	case <-ctx.Done():
		fmt.Println("worker is canceled id: ", id, " ", ctx.Err())

	}

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	duratios := []int{1, 2, 3}
	done := make(chan int, len(duratios))

	for id, duration := range duratios {
		wg.Add(1)
		go doWorkCorrect(ctx, id+1, duration, done, &wg)
	}
	firstFinish := <-done
	fmt.Println("fist finish: ", firstFinish)

	cancel()

	wg.Wait()

}
