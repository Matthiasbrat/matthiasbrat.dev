---
title: "Go Basics Part 1: Getting Started"
description: "Learn Go programming fundamentals - variables, types, and control structures"
date: 2025-01-10
---

Welcome to this series on Go programming! In this three-part series, we'll cover Go fundamentals, from basic syntax to building real applications.

## What is Go?

Go (also called Golang) is a statically typed, compiled programming language designed at Google. It's known for:

- Simple, clean syntax
- Fast compilation
- Built-in concurrency
- Excellent standard library
- Strong tooling

::: aside
**History**: Go was created in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson at Google. It was publicly announced in November 2009.
:::

## Installing Go

### On macOS

```bash
brew install go
```

### On Ubuntu/Debian

```bash
sudo apt update
sudo apt install golang-go
```

###On Windows

Download the installer from [golang.org](https://golang.org/dl/)

> [!TIP]
> After installation, verify with `go version`

## Your First Go Program

Create a file called `hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Run it with:

```bash
go run hello.go
```

> [!NOTE]
> Every Go program starts with a `package` declaration. The `main` package is special - it defines a standalone executable.

## Variables and Types

### Variable Declaration

```go
// Method 1: var keyword with type
var name string = "Alice"
var age int = 30

// Method 2: var keyword with type inference
var city = "New York"

// Method 3: short declaration (inside functions only)
country := "USA"
```

::: aside
The `:=` operator is Go's **short variable declaration**. It's concise and commonly used, but only works inside functions.
:::

### Basic Types

```go
// Strings
var message string = "Hello"

// Integers
var count int = 42
var smallNumber int8 = 127
var bigNumber int64 = 9223372036854775807

// Floating point
var price float64 = 19.99
var discount float32 = 0.15

// Boolean
var isActive bool = true

// Constants
const Pi = 3.14159
const AppName = "MyApp"
```

> [!IMPORTANT]
> Go is statically typed. Once a variable is declared with a type, it cannot change.

## Control Structures

### If Statements

```go
age := 18

if age >= 18 {
    fmt.Println("Adult")
} else if age >= 13 {
    fmt.Println("Teenager")
} else {
    fmt.Println("Child")
}

// If with initialization
if err := doSomething(); err != nil {
    fmt.Println("Error:", err)
}
```

### For Loops

Go only has `for` loops (no `while` or `do-while`):

```go
// Traditional for loop
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// While-like loop
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}

// Infinite loop
for {
    // Use break to exit
    if someCondition {
        break
    }
}

// Range over slice
numbers := []int{1, 2, 3, 4, 5}
for index, value := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}
```

### Switch Statements

```go
day := "Monday"

switch day {
case "Monday":
    fmt.Println("Start of the week")
case "Friday":
    fmt.Println("TGIF!")
case "Saturday", "Sunday":
    fmt.Println("Weekend!")
default:
    fmt.Println("Midweek")
}

// Switch without expression
hour := 14
switch {
case hour < 12:
    fmt.Println("Morning")
case hour < 18:
    fmt.Println("Afternoon")
default:
    fmt.Println("Evening")
}
```

> [!TIP]
> Unlike C or Java, Go's switch cases don't fall through by default. No need for `break` statements!

## Functions

```go
// Basic function
func greet(name string) {
    fmt.Println("Hello,", name)
}

// Function with return value
func add(a int, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Named return values
func calculate(a, b int) (sum int, product int) {
    sum = a + b
    product = a * b
    return // naked return
}
```

## Arrays and Slices

### Arrays

```go
// Fixed-size array
var numbers [5]int
numbers[0] = 1
numbers[1] = 2

// Array literal
primes := [5]int{2, 3, 5, 7, 11}
```

### Slices

```go
// Create a slice
fruits := []string{"apple", "banana", "cherry"}

// Append to slice
fruits = append(fruits, "date")

// Slice of slice
subset := fruits[1:3] // banana, cherry

// Make a slice with capacity
numbers := make([]int, 5, 10) // length 5, capacity 10
```

> [!WARNING]
> Slices are references to underlying arrays. Modifying a slice can affect other slices sharing the same array.

## Maps

```go
// Create a map
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}

// Add/update
ages["Charlie"] = 35

// Get value
age := ages["Alice"]

// Check if key exists
age, exists := ages["David"]
if exists {
    fmt.Println("Age:", age)
}

// Delete key
delete(ages, "Bob")

// Iterate over map
for name, age := range ages {
    fmt.Printf("%s is %d years old\n", name, age)
}
```

## Task List

- [x] Install Go
- [x] Learn variables and types
- [x] Understand control structures
- [x] Write functions
- [x] Use arrays, slices, and maps
- [ ] Learn about pointers (Part 2)
- [ ] Understand structs and methods (Part 2)
- [ ] Master concurrency (Part 3)

## Next Steps

In [Part 2](/blog/go-basics-part-2), we'll dive into:
- Pointers
- Structs and methods
- Interfaces
- Error handling

Keep practicing with small programs to solidify these concepts!
