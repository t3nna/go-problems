package main

import (
	"context"
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
	resCh := make(chan string, 1)

	for _, addr := range addresses {
		go func(addr string) {
			res, err := getter.Get(ctx, addr, key)

			if err != nil {
				errCh <- err
			} else {
				select {
				case resCh <- res:
				default:

				}
			}

		}(addr)

	}

	errCount := 0

	for {
		select {
		case val := <-errCh:
			errCount++
			if errCount == len(addresses) {
				return "", val
			}
		case <-ctx.Done():
			return "", context.Canceled
		case res := <-resCh:
			return res, nil
		}
	}

}
