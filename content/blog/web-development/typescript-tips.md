---
title: "TypeScript Tips for Better Code"
description: "Advanced TypeScript patterns that will level up your code"
date: 2025-02-10
---

TypeScript's type system is incredibly powerful. Here are some tips to use it effectively.

## Utility Types

Built-in types that save time:

```typescript
interface User {
    id: number;
    name: string;
    email: string;
    password: string;
}

// Omit sensitive fields
type PublicUser = Omit<User, 'password'>;

// Make all fields optional
type PartialUser = Partial<User>;

// Make all fields required
type RequiredUser = Required<User>;

// Pick specific fields
type UserCredentials = Pick<User, 'email' | 'password'>;
```

## Discriminated Unions

Type-safe handling of different cases:

```typescript
type Result<T> =
    | { success: true; data: T }
    | { success: false; error: string };

function handleResult(result: Result<User>) {
    if (result.success) {
        console.log(result.data.name);  // TypeScript knows data exists
    } else {
        console.log(result.error);  // TypeScript knows error exists
    }
}
```

> [!TIP]
> Use discriminated unions instead of optional fields when possible.

## Template Literal Types

```typescript
type EventName = 'click' | 'focus' | 'blur';
type Handler = `on${Capitalize<EventName>}`;
// Handler = 'onClick' | 'onFocus' | 'onBlur'
```

## Const Assertions

```typescript
const config = {
    api: 'https://api.example.com',
    timeout: 5000,
} as const;

// config.api is type 'https://api.example.com', not string
```

> [!NOTE]
> `as const` makes the object deeply readonly with literal types.
