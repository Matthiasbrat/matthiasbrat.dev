---
title: "Variables and Types in Go"
description: "Understanding Go's type system and variable declarations"
date: 2025-01-20
---

Go is statically typed, meaning variable types are known at compile time. Let's explore how to work with variables and types.

## Variable Declaration

```go
// Explicit type
var name string = "Alice"
var age int = 30

// Type inference
var city = "New York"

// Short declaration (functions only)
country := "USA"
```

> [!NOTE]
> The `:=` operator is only available inside functions.

## Basic Types

```go
// Strings
var message string = "Hello"

// Integers
var count int = 42
var bigNumber int64 = 9223372036854775807

// Floating point
var price float64 = 19.99

// Boolean
var isActive bool = true

// Constants
const Pi = 3.14159
```

## Zero Values

Uninitialized variables get zero values:

| Type | Zero Value |
|------|------------|
| int | 0 |
| float64 | 0.0 |
| string | "" |
| bool | false |
| pointer | nil |

## Type Conversions

Go requires explicit type conversions:

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)
```

> [!WARNING]
> Be careful with type conversions - you might lose precision or overflow.
