---
title: "Getting Started with Go"
description: "Install Go and write your first program in under 10 minutes"
date: 2025-01-15
---

Go is a statically typed, compiled programming language designed at Google. It's known for simplicity, performance, and excellent concurrency support.

## Installing Go

### macOS

```bash
brew install go
```

### Linux

```bash
sudo apt update && sudo apt install golang-go
```

### Windows

Download the installer from [golang.org](https://golang.org/dl/).

> [!TIP]
> Verify installation with `go version`

## Your First Program

Create `hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Run it:

```bash
go run hello.go
```

## What's Next

In the next post, we'll explore variables, types, and control structures.
