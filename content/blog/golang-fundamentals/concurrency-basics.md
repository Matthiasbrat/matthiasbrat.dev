---
title: "Concurrency in Go"
description: "Goroutines and channels - Go's powerful concurrency primitives"
date: 2025-01-25
---

Go's concurrency model is one of its most powerful features. Built around goroutines and channels, it makes concurrent programming surprisingly approachable.

## Goroutines

A goroutine is a lightweight thread managed by the Go runtime:

```go
func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    go sayHello("Alice")  // Runs concurrently
    go sayHello("Bob")    // Also concurrent

    time.Sleep(100 * time.Millisecond)
}
```

> [!TIP]
> Goroutines are extremely cheap - you can easily run thousands of them.

## Channels

Channels are typed conduits for communication between goroutines:

```go
func main() {
    messages := make(chan string)

    go func() {
        messages <- "Hello from goroutine!"
    }()

    msg := <-messages
    fmt.Println(msg)
}
```

## Buffered Channels

```go
ch := make(chan int, 3)  // Buffer size 3

ch <- 1
ch <- 2
ch <- 3
// ch <- 4 would block until someone reads
```

## Select Statement

Handle multiple channels:

```go
select {
case msg := <-ch1:
    fmt.Println("From ch1:", msg)
case msg := <-ch2:
    fmt.Println("From ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout!")
}
```

> [!IMPORTANT]
> Always ensure goroutines can exit. Leaked goroutines are a common source of bugs.
