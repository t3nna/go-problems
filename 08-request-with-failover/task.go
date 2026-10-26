package main

import (
	"context"
	"errors"
	"fmt"
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
	errCh := make(chan error, len(addresses))
	resCh := make(chan string, len(addresses))

	errCount := 0

	for i := 0; i < len(addresses); i++ {
		v := addresses[i]
		innerCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()

		go func() {
			str, err := client.Get(innerCtx, v)
			if err != nil {
				errCh <- err
			} else {
				resCh <- str
			}
		}()

		select {
		case <-innerCtx.Done():
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
			}
		case res := <-resCh:
			return res, nil
		case <-errCh:
			errCount++
		}
		fmt.Println("innner", v)

	}

	if errCount == len(addresses) {
		return "", ErrRequestsFailed
	}

	for {
		select {
		case v := <-resCh:
			return v, nil
		case <-errCh:
			errCount++
			if errCount == len(addresses) {
				return "", ErrRequestsFailed
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
