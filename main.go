package main

import (
	"sync"
)

type stateManager struct {
	fetched   map[string]bool
	waitGroup *sync.WaitGroup
	rwMutex   *sync.RWMutex
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func main() {
	sm := stateManager{
		fetched:   make(map[string]bool),
		waitGroup: new(sync.WaitGroup),
		rwMutex:   &sync.RWMutex{},
	}
	Crawl("https://golang.org/", 3, fetcher, sm)
	sm.waitGroup.Wait()
}
