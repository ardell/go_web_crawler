package main

import (
	"fmt"
	"runtime"
)

func isAlreadyFetched(url string, sm stateManager) bool {
	sm.rwMutex.RLock()
	defer sm.rwMutex.RUnlock()

	var _, ok = sm.fetched[url]
	return ok
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, sm stateManager) {
	sm.waitGroup.Add(1)

	go func() {
		// Tell the WaitGroup we're finished whenever we leave this scope
		defer sm.waitGroup.Done()

		// Don't process further than the specified depth
		if depth <= 0 {
			return
		}

		// Safely return if we've already fetched the given URL
		if isAlreadyFetched(url, sm) {
			fmt.Printf("Skipping url %s (already fetched)\n", url)
			return
		}

		// Handle failure to fetch
		fmt.Printf("Fetching url %s at depth %v\n", url, depth)
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("  found body %s\n", body)
		}

		// Add the fetched URL to SafeFetcher
		sm.rwMutex.Lock()
		sm.fetched[url] = true
		sm.rwMutex.Unlock()

		// Continue crawling
		for _, u := range urls {
			fmt.Printf("  dispatching url %s\n", u)
			Crawl(u, depth-1, fetcher, sm)

			// Yield so other threads can make progress
			runtime.Gosched()
		}
	}()

	return
}
