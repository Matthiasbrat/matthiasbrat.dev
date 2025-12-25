---
title: "Go Basics Part 3: Concurrency"
description: "Master Go's concurrency primitives - goroutines, channels, and patterns"
date: 2025-01-14
---

This is the final part of our Go basics series. After covering [fundamentals](/blog/go-basics-part-1) and [types](/blog/go-basics-part-2), let's explore Go's killer feature: **concurrency**.

## Goroutines

Goroutines are lightweight threads managed by the Go runtime.

```go
func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    // Launch goroutine
    go sayHello()

    // Without this, program exits before goroutine runs
    time.Sleep(time.Second)
}
```

::: aside
**Lightweight**: You can run thousands of goroutines. They start with just 2KB of stack space and grow as needed.
:::

### Anonymous Goroutines

```go
func main() {
    go func() {
        fmt.Println("Anonymous goroutine")
    }()

    time.Sleep(time.Second)
}
```

> [!CAUTION]
> Using `time.Sleep` to wait for goroutines is bad practice. Use channels or `sync.WaitGroup` instead!

## Channels

Channels allow goroutines to communicate safely.

```go
func main() {
    // Create channel
    messages := make(chan string)

    // Send in goroutine
    go func() {
        messages <- "Hello"
    }()

    // Receive (blocks until value available)
    msg := <-messages
    fmt.Println(msg)
}
```

### Buffered Channels

```go
// Unbuffered (blocks on send until received)
ch1 := make(chan int)

// Buffered (can hold 3 values before blocking)
ch2 := make(chan int, 3)

ch2 <- 1
ch2 <- 2
ch2 <- 3
// Next send would block until something is received
```

### Channel Directions

```go
// Send-only channel
func sendOnly(ch chan<- string) {
    ch <- "message"
}

// Receive-only channel
func receiveOnly(ch <-chan string) {
    msg := <-ch
    fmt.Println(msg)
}

func main() {
    ch := make(chan string)
    go sendOnly(ch)
    receiveOnly(ch)
}
```

## Select Statement

The `select` statement lets you wait on multiple channel operations.

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "one"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "two"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        }
    }
}
```

### Select with Timeout

```go
select {
case msg := <-messages:
    fmt.Println("Received:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
}
```

### Select with Default

```go
select {
case msg := <-messages:
    fmt.Println("Received:", msg)
default:
    fmt.Println("No message received")
}
```

## Common Patterns

### Worker Pool

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for a := 1; a <= 5; a++ {
        <-results
    }
}
```

> [!TIP]
> Close channels when done sending to signal workers to exit.

### Rate Limiting

```go
func main() {
    requests := make(chan int, 5)
    for i := 1; i <= 5; i++ {
        requests <- i
    }
    close(requests)

    limiter := time.Tick(200 * time.Millisecond)

    for req := range requests {
        <-limiter // wait for tick
        fmt.Println("Request", req, time.Now())
    }
}
```

### Pipeline

```go
func generateNumbers(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    // Set up pipeline
    numbers := generateNumbers(2, 3, 4, 5)
    squares := square(numbers)

    // Consume output
    for result := range squares {
        fmt.Println(result)
    }
}
```

## sync.WaitGroup

Wait for multiple goroutines to complete.

```go
import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
    fmt.Println("All workers done")
}
```

## sync.Mutex

Protect shared state with mutexes.

```go
import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    c.v[key]++
    c.mu.Unlock()
}

func (c *SafeCounter) Value(key string) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.v[key]
}

func main() {
    counter := SafeCounter{v: make(map[string]int)}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc("count")
        }()
    }

    wg.Wait()
    fmt.Println(counter.Value("count")) // 1000
}
```

> [!IMPORTANT]
> Always use `defer mu.Unlock()` to ensure the mutex is released even if a panic occurs.

## Real-World Example: Web Scraper

```go
package main

import (
    "fmt"
    "net/http"
    "sync"
)

func fetchURL(url string, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done()

    resp, err := http.Get(url)
    if err != nil {
        results <- fmt.Sprintf("Error fetching %s: %v", url, err)
        return
    }
    defer resp.Body.Close()

    results <- fmt.Sprintf("%s: Status %d", url, resp.StatusCode)
}

func main() {
    urls := []string{
        "https://golang.org",
        "https://github.com",
        "https://stackoverflow.com",
    }

    var wg sync.WaitGroup
    results := make(chan string, len(urls))

    for _, url := range urls {
        wg.Add(1)
        go fetchURL(url, &wg, results)
    }

    // Close results when all goroutines done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Print results as they arrive
    for result := range results {
        fmt.Println(result)
    }
}
```

## Concurrency vs Parallelism

::: aside
**Concurrency** is about *dealing* with multiple things at once.
**Parallelism** is about *doing* multiple things at once.

Go provides concurrency. The runtime handles parallelism across CPU cores automatically.
:::

> [!NOTE]
> Set `GOMAXPROCS` to control how many OS threads can execute Go code simultaneously. Default is the number of CPU cores.

## Best Practices

- [ ] Don't communicate by sharing memory; share memory by communicating
- [ ] Use channels for communication between goroutines
- [ ] Close channels to signal completion
- [ ] Always handle receive from closed channels
- [ ] Avoid goroutine leaks (always ensure goroutines can exit)
- [ ] Use `context.Context` for cancellation in real applications

## Series Wrap-Up

Congratulations! You've completed the Go Basics series:

| Part | Topics | Link |
|------|--------|------|
| **Part 1** | Variables, control flow, functions | [Read →](/blog/go-basics-part-1) |
| **Part 2** | Pointers, structs, interfaces | [Read →](/blog/go-basics-part-2) |
| **Part 3** | Goroutines, channels, concurrency | ✅ Current |

## Next Steps

- Build a real project (web server, CLI tool, API)
- Read [Effective Go](https://golang.org/doc/effective_go)
- Explore the standard library
- Join the Go community

Happy coding! 🚀
