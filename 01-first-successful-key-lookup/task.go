package main

import (
	"context"
	"sync"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if len(addresses) == 0 {
		return "", nil
	}

	resCh, errCh := make(chan string, len(addresses)), make(chan error, len(addresses))

	var wg sync.WaitGroup

	wg.Add(len(addresses))

	for _, v := range addresses {
		go func(addr string) {
			res, err := getter.Get(ctx, addr, key)
			if err != nil {
				errCh <- err
			} else {
				resCh <- res
			}

			wg.Done()
		}(v)

	}

	go func() {
		wg.Wait()
		close(resCh)
		close(errCh)
	}()

	errCounter := 0
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case v := <-errCh:
			errCounter++
			if errCounter == len(addresses) {
				return "", v
			}
		case v := <-resCh:
			return v, nil

		}
	}
}
