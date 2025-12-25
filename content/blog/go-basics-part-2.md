---
title: "Go Basics Part 2: Structs and Interfaces"
description: "Learn about Go's type system - pointers, structs, methods, and interfaces"
date: 2025-01-12
---

In [Part 1](/blog/go-basics-part-1), we covered Go fundamentals. Now let's explore Go's type system and object-oriented patterns.

## Pointers

Go has pointers, but no pointer arithmetic (unlike C).

```go
func main() {
    x := 10
    p := &x // pointer to x

    fmt.Println(*p) // dereference: 10
    *p = 20         // modify via pointer
    fmt.Println(x)  // 20
}
```

::: aside
**Memory Safety**: Go's pointers are safer than C pointers because they don't support arithmetic operations that could access invalid memory.
:::

### Pointers with Functions

```go
// Without pointer - value is copied
func doubleValue(x int) {
    x = x * 2 // modifies copy, not original
}

// With pointer - can modify original
func doublePointer(x *int) {
    *x = *x * 2 // modifies original
}

func main() {
    num := 5
    doubleValue(num)
    fmt.Println(num) // Still 5

    doublePointer(&num)
    fmt.Println(num) // Now 10
}
```

> [!TIP]
> Use pointers to avoid copying large structs and when you need to modify the original value.

## Structs

Structs are Go's way of creating custom types.

```go
type Person struct {
    Name string
    Age  int
    Email string
}

// Create struct instances
alice := Person{
    Name:  "Alice",
    Age:   30,
    Email: "alice@example.com",
}

// Shorter syntax (order matters)
bob := Person{"Bob", 25, "bob@example.com"}

// Access fields
fmt.Println(alice.Name)
alice.Age = 31
```

### Embedded Structs

```go
type Address struct {
    Street string
    City   string
}

type Employee struct {
    Person  // embedding
    Address // embedding
    Salary  int
}

emp := Employee{
    Person: Person{Name: "Alice", Age: 30},
    Address: Address{Street: "123 Main St", City: "NYC"},
    Salary: 100000,
}

// Access embedded fields directly
fmt.Println(emp.Name)   // from Person
fmt.Println(emp.City)   // from Address
fmt.Println(emp.Salary) // from Employee
```

## Methods

Methods are functions with a receiver.

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method with pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}

    area := rect.Area()
    fmt.Println("Area:", area)

    rect.Scale(2)
    fmt.Println("New dimensions:", rect.Width, rect.Height)
}
```

> [!IMPORTANT]
> Use pointer receivers when:
> - You need to modify the receiver
> - The struct is large (avoid copying)
> - For consistency (if some methods need pointers, use pointers for all)

## Interfaces

Interfaces define behavior. Types implement interfaces *implicitly*.

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// Rectangle already implements Shape (has Area and Perimeter)

func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
    circle := Circle{Radius: 5}
    rect := Rectangle{Width: 10, Height: 5}

    printShapeInfo(circle)
    printShapeInfo(rect)
}
```

::: aside
**Duck Typing**: If it walks like a duck and quacks like a duck, it's a duck. In Go, if a type has the required methods, it implements the interface automatically.
:::

### Empty Interface

```go
// interface{} accepts any type
func printAnything(value interface{}) {
    fmt.Println(value)
}

printAnything("Hello")
printAnything(42)
printAnything([]int{1, 2, 3})
```

### Type Assertions

```go
func describeValue(i interface{}) {
    // Type assertion
    if s, ok := i.(string); ok {
        fmt.Printf("String: %s\n", s)
    } else if n, ok := i.(int); ok {
        fmt.Printf("Integer: %d\n", n)
    }

    // Type switch
    switch v := i.(type) {
    case string:
        fmt.Printf("String of length %d\n", len(v))
    case int:
        fmt.Printf("Integer: %d\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

## Error Handling

Go uses explicit error returns instead of exceptions.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

### Custom Errors

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field:   "age",
            Message: "must be positive",
        }
    }
    if age > 150 {
        return &ValidationError{
            Field:   "age",
            Message: "unrealistic value",
        }
    }
    return nil
}
```

> [!WARNING]
> Always check for errors! Ignoring errors is a common source of bugs in Go programs.

## Comparison Table

| Feature | Value Receiver | Pointer Receiver |
|---------|---------------|------------------|
| Modifies receiver | ❌ No | ✅ Yes |
| Works with nil | ❌ No | ✅ Yes |
| Copies struct | ✅ Yes | ❌ No (just pointer) |
| Use when | Small structs, read-only | Large structs, modifications needed |

## Coming Up

In [Part 3](/blog/go-basics-part-3), we'll explore:
- Goroutines and channels
- Concurrency patterns
- The `select` statement
- Real-world examples

The final part will bring everything together!
