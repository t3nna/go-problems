package main

import (
	"sync"
)

func Merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	res := make(chan int)

	for _, v := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			for v := range c {
				res <- v
			}
			wg.Done()
		}(v)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
