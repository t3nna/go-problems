package main

import (
	"fmt"
	"sync"
)

func Merge(channels ...<-chan int) <-chan int {
	res := make(chan int)

	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				fmt.Println(val)
				res <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
