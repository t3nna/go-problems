package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var (
	mu sync.RWMutex
	c  = map[string]bool{}
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) ([]string, error) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return nil, nil
	}

	mu.RLock()
	if c[url] {
		mu.RUnlock()
		return nil, nil
	}
	mu.RUnlock()

	body, urls, err := fetcher.Fetch(url)

	if err != nil {
		return nil, err
	}
	mu.Lock()
	c[url] = true
	mu.Unlock()

	result := []string{body}

	rCh := make(chan string, 1)
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(u)
			if res, err := Crawl(u, depth-1, fetcher); err == nil {
				if res == nil {
					return
				}
				for _, b := range res {
					fmt.Println(b)
					rCh <- b
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(rCh)
	}()

	for v := range rCh {
		result = append(result, v)
	}

	return result, nil
}
