---
title: "Endpoints"
description: "Complete list of API endpoints"
order: 3
---

Complete reference for all available API endpoints.

## Users

### List Users

```http
GET /v1/users
```

**Parameters:**

| Name | Type | Description |
|------|------|-------------|
| `limit` | integer | Max results (default: 20, max: 100) |
| `offset` | integer | Pagination offset |
| `status` | string | Filter by status: `active`, `inactive` |

**Response:**

```json
{
  "data": [
    {
      "id": "usr_abc123",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "total": 42,
    "limit": 20,
    "offset": 0
  }
}
```

### Get User

```http
GET /v1/users/:id
```

### Create User

```http
POST /v1/users
```

**Request Body:**

```json
{
  "email": "newuser@example.com",
  "name": "Jane Smith",
  "role": "member"
}
```

### Update User

```http
PATCH /v1/users/:id
```

### Delete User

```http
DELETE /v1/users/:id
```

> [!CAUTION]
> Deleting a user is permanent and cannot be undone.

## Projects

### List Projects

```http
GET /v1/projects
```

### Create Project

```http
POST /v1/projects
```

```json
{
  "name": "My Project",
  "description": "A great project"
}
```

## Webhooks

### List Webhooks

```http
GET /v1/webhooks
```

### Create Webhook

```http
POST /v1/webhooks
```

```json
{
  "url": "https://yourapp.com/webhook",
  "events": ["user.created", "user.deleted"]
}
```

> [!TIP]
> Use webhook signatures to verify that requests are from us.
