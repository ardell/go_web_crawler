# Exercise: Web Crawler

From [A Tour of Go](https://go.dev/tour/concurrency/10), concurrency exercise #10.

## To run

```
$> go run *.go
```

## Explanation

1. In `main.go` we create a `stateManager` struct, which holds 3 things:
  a. a `fetched` map, so we can check whether we've already fetched a URL in
  `O(1)` time,
  b. a `waitGroup` that will `Add(1)` when we start fetching a url and `Done`
  when we finish fetching a url, and
  c. a `rwMutex` so we can safely read to and write from the `fetched` array
  from many goroutines concurrently

2. In `crawl.go` we implement an `isAlreadyFetched` helper function that uses
   our `stateManager.rwMutex` to check whether a url has already been fetched
   and return a boolean.

3. Also in `crawl.go`, we:
  a. `Add(1)` to our `stateManager.waitGroup` to tell the `WaitGroup` that we
  are waiting to fetch an item,
  b. implement a `goroutine` using an IIFE (Immediately Invoked Function
  Expression) -- `go func() { ... }()` that does the following...
    1. tell the wait group that we're finished whenever we return from the IIFE's
    scope,
    2. return if we've reached our maximum fetching depth,
    3. skip processing the node if we've already fetched the url,
    4. handle fetching failures,
    5. safely mark the url as fetched in `stateManager.fetched` using our
       `stateManager.rwMutex` (lock for _writing_ this time), and
    6. recurse and yield to other goroutines.

## Future improvements

1. Allow `main.go` to receive the depth as a command-line argument.
2. Refactor the `Crawl` function's IIFE. It has too many responsibilities, and
   mixes levels of abstraction: orchestration of concurrency, logging, error
   handling, recursion.
