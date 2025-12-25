---
title: "REST API Design Best Practices"
description: "Building intuitive, consistent, and scalable REST APIs"
date: 2025-02-05
---

Good API design makes your service a joy to use. Here are battle-tested practices for REST APIs.

## Resource Naming

Use nouns, not verbs. Use plural forms:

```
GET    /users          # List users
GET    /users/123      # Get user 123
POST   /users          # Create user
PUT    /users/123      # Update user 123
DELETE /users/123      # Delete user 123
```

> [!WARNING]
> Avoid: `/getUsers`, `/createUser`, `/user/123`

## HTTP Status Codes

Use them correctly:

| Code | Meaning | When to Use |
|------|---------|-------------|
| 200 | OK | Successful GET/PUT |
| 201 | Created | Successful POST |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Validation error |
| 401 | Unauthorized | Missing/invalid auth |
| 403 | Forbidden | No permission |
| 404 | Not Found | Resource doesn't exist |
| 500 | Server Error | Bug in your code |

## Pagination

Always paginate list endpoints:

```json
{
    "data": [...],
    "pagination": {
        "page": 1,
        "per_page": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

## Error Responses

Be consistent and helpful:

```json
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid email format",
        "field": "email"
    }
}
```

> [!TIP]
> Include enough detail to help clients fix the issue, but never expose internal errors.
