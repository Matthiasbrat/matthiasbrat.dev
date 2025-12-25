---
title: "Complete Markdown Syntax Guide"
description: "A comprehensive guide showcasing all available markdown syntax and features on this site"
date: 2025-01-15
---

This post demonstrates every markdown feature available on this site. Use it as a reference when writing content.

## Headings

# Heading 1

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

###### Heading 6

## Text Formatting

This is **bold text** and this is **also bold**.

This is _italic text_ and this is _also italic_.

This is **_bold and italic_** text.

This is ~~strikethrough~~ text.

## Links and Images

[This is a link](https://example.com)

[This is a link with a title](https://example.com "Example Website")

![Alt text for an image](https://via.placeholder.com/600x300)

## Lists

### Unordered List

- Item 1
- Item 2
  - Nested item 2.1
  - Nested item 2.2
- Item 3

### Ordered List

1. First item
2. Second item
   1. Nested item 2.1
   2. Nested item 2.2
3. Third item

### Task List

- [x] Completed task
- [x] Another completed task
- [ ] Incomplete task
- [ ] Another incomplete task

## Blockquotes

> This is a blockquote.
> It can span multiple lines.
>
> And even include multiple paragraphs.

## Code

### Inline Code

Use `var x = 10;` for inline code.

### Code Blocks

#### JavaScript

```javascript
function greet(name) {
  console.log(`Hello, ${name}!`);
  return `Welcome, ${name}`;
}

const user = "Alice";
greet(user);
```

#### Go

```go
package main

import "fmt"

func main() {
    message := "Hello, World!"
    fmt.Println(message)

    for i := 0; i < 5; i++ {
        fmt.Printf("Count: %d\n", i)
    }
}
```

#### Python

```python
def calculate_fibonacci(n):
    if n <= 1:
        return n
    return calculate_fibonacci(n-1) + calculate_fibonacci(n-2)

# Generate first 10 Fibonacci numbers
for i in range(10):
    print(f"F({i}) = {calculate_fibonacci(i)}")
```

#### TypeScript

```typescript
interface User {
  id: number;
  name: string;
  email: string;
}

class UserService {
  private users: User[] = [];

  addUser(user: User): void {
    this.users.push(user);
  }

  getUserById(id: number): User | undefined {
    return this.users.find((u) => u.id === id);
  }
}
```

## Tables

| Feature             | Supported | Notes                    |
| ------------------- | --------- | ------------------------ |
| Markdown            | ✅        | Full CommonMark          |
| Syntax Highlighting | ✅        | 180+ languages           |
| Tables              | ✅        | GitHub Flavored Markdown |
| Task Lists          | ✅        | Interactive checkboxes   |
| Callouts            | ✅        | 5 types available        |

## Horizontal Rules

---

---

---

## Callouts

> [!NOTE]
> This is a note callout. Use it for informational content that readers should be aware of.

> [!TIP]
> This is a tip callout. Share helpful suggestions and best practices here.

> [!IMPORTANT]
> This is an important callout. Highlight critical information that readers must not miss.

> [!WARNING]
> This is a warning callout. Alert readers about potential issues or pitfalls.

> [!CAUTION]
> This is a caution callout. Use for critical warnings about dangerous operations or irreversible actions.

### Foldable Callouts

> [!TIP]-
> This is a foldable callout! It starts collapsed.
>
> Click the title to expand and see more content. Perfect for optional details or advanced topics.

## Asides

::: aside
**Did you know?** Asides provide additional context without interrupting the main flow. On desktop, they appear to the right of the content. On mobile, they're displayed inline.
:::

This is regular paragraph text that flows around the aside on desktop. The aside floats to the right and doesn't interrupt the reading flow.

::: aside
Asides are perfect for:

- Biographical information
- Historical context
- Related topics
- Fun facts
- Additional resources
  :::

You can have multiple asides in your content, and they will stack vertically on the right side.

## Embedded Content

### YouTube Videos

![Introduction to Go Programming](https://www.youtube.com/watch?v=YS4e4q9oBaU)

## Advanced Formatting

### Nested Lists with Code

1. First step: Install dependencies

   ```bash
   npm install -D tailwindcss
   ```

2. Second step: Configure Tailwind

   - Create `tailwind.config.js`
   - Add to PostCSS config

   ```javascript
   module.exports = {
     content: ["./src/**/*.{html,js}"],
     theme: {
       extend: {},
     },
   };
   ```

3. Third step: Build your CSS

### Combining Features

> [!TIP]
> You can combine different markdown features for rich content.
>
> For example, this callout includes:
>
> - **Bold text**
> - _Italic text_
> - `inline code`
> - [Links](https://example.com)
>
> ```go
> // And even code blocks!
> func main() {
>     fmt.Println("Hello from inside a callout!")
> }
> ```

## Summary

This guide covered all available markdown syntax:

- ✅ Headings (H1-H6)
- ✅ Text formatting (bold, italic, strikethrough)
- ✅ Links and images
- ✅ Lists (ordered, unordered, task lists)
- ✅ Blockquotes
- ✅ Inline and block code
- ✅ Syntax highlighting (180+ languages)
- ✅ Tables
- ✅ Horizontal rules
- ✅ Callouts (NOTE, TIP, IMPORTANT, WARNING, CAUTION)
- ✅ Foldable callouts
- ✅ Asides
- ✅ YouTube embeds

Happy writing!
