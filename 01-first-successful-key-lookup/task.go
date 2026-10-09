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

	if len(addresses) == 0 {
		return "", nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(addresses))
	resCh := make(chan string, len(addresses))
	var wg sync.WaitGroup

	wg.Add(len(addresses))

	for _, v := range addresses {

		go func(s string) {
			res, err := getter.Get(ctx, s, key)
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
		close(errCh)
		close(resCh)
	}()

	errCount := 0

	for {
		select {
		case v := <-errCh:
			errCount++
			if errCount == len(addresses) {
				return "", v
			}
		case <-ctx.Done():
			return "", ctx.Err()
		case v := <-resCh:
			return v, nil

		}
	}

}
