---
title: "SDKs & Libraries"
description: "Official SDKs and community libraries"
order: 5
---

Use our official SDKs for easier integration.

## Official SDKs

### JavaScript / TypeScript

```bash
npm install @example/sdk
```

```javascript
import { Client } from '@example/sdk';

const client = new Client({
  apiKey: process.env.API_KEY
});

// List users
const users = await client.users.list({ limit: 10 });

// Create a user
const newUser = await client.users.create({
  email: 'user@example.com',
  name: 'John Doe'
});
```

### Python

```bash
pip install example-sdk
```

```python
from example import Client

client = Client(api_key="sk_live_...")

# List users
users = client.users.list(limit=10)

# Create a user
new_user = client.users.create(
    email="user@example.com",
    name="John Doe"
)
```

### Go

```bash
go get github.com/example/sdk-go
```

```go
package main

import (
    "github.com/example/sdk-go"
)

func main() {
    client := sdk.NewClient("sk_live_...")

    // List users
    users, err := client.Users.List(&sdk.ListParams{
        Limit: 10,
    })

    // Create a user
    user, err := client.Users.Create(&sdk.CreateUserParams{
        Email: "user@example.com",
        Name:  "John Doe",
    })
}
```

### Ruby

```bash
gem install example-sdk
```

```ruby
require 'example'

client = Example::Client.new(api_key: 'sk_live_...')

# List users
users = client.users.list(limit: 10)

# Create a user
new_user = client.users.create(
  email: 'user@example.com',
  name: 'John Doe'
)
```

## Community Libraries

These libraries are maintained by the community:

| Language | Library | Maintainer |
|----------|---------|------------|
| PHP | [example-php](https://github.com/...) | @contributor |
| Rust | [example-rs](https://github.com/...) | @rustacean |
| Elixir | [example_ex](https://github.com/...) | @alchemist |

> [!NOTE]
> Community libraries are not officially supported. Use at your own risk.

## Direct HTTP

If there's no SDK for your language, you can use HTTP directly:

```bash
curl -X POST "https://api.example.com/v1/users" \
     -H "Authorization: Bearer sk_live_..." \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "name": "John Doe"}'
```
