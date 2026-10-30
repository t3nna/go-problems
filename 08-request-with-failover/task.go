package main

import (
	"context"
	"errors"
	"time"
)

type Client interface {
	Get(ctx context.Context, address string) (string, error)
}

var ErrRequestsFailed = errors.New("task execution failed")

// RequestWithFailover attempts to request a data from available addresses:
// 1. If error, immediately try the address without waiting
// 2. If an address doesn't respond within 500ms, try the next but keep the original request running
// 3. Return the first successful response, or all ErrRequestsFailed if all nodes fail
// 4. Properly handle context cancellation throughout the process
func RequestWithFailover(ctx context.Context, client Client, addresses []string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	resCh := make(chan string, len(addresses))
	errCh := make(chan error, len(addresses))

	errCounter := 0
	for _, addr := range addresses {
		ctxInnr, cancel := context.WithTimeout(ctx, 500*time.Millisecond)

		go func(a string) {
			res, err := client.Get(ctx, a)
			if err != nil {
				errCh <- err
			} else {
				resCh <- res
			}
		}(addr)

		defer cancel()
		select {
		case <-ctxInnr.Done():
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
			}

		case <-errCh:
			errCounter++
			if errCounter == len(addresses) {
				return "", ErrRequestsFailed
			}
		case v := <-resCh:
			return v, nil

		}

	}

	for {
		select {
		case <-errCh:
			errCounter++
			if errCounter == len(addresses) {
				return "", ErrRequestsFailed
			}
		case <-ctx.Done():
			return "", ctx.Err()

		case v := <-resCh:
			return v, nil

		}

	}
}
