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

	resCh := make(chan string, 1)
	errCh := make(chan error, len(addresses))

	for _, val := range addresses {
		go func(s string) {
			res, err := getter.Get(ctx, s, key)
			if err != nil {
				errCh <- err
				return
			}
			select {

			case resCh <- res:
			default:
			}
		}(val)
	}

	errCounter := 0
	for {
		select {
		case <-ctx.Done():
			return "", context.Canceled

		case v := <-errCh:
			errCounter++
			if errCounter == len(addresses) {
				return "", v
			}
		case r := <-resCh:
			return r, nil
		}

	}

}
