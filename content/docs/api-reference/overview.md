---
title: "Overview"
description: "Introduction to the API and core concepts"
order: 1
---

This section provides complete documentation for our API.

## Base URL

All API requests should be made to:

```
https://api.example.com/v1
```

## Authentication

The API uses Bearer token authentication:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
     https://api.example.com/v1/users
```

## Rate Limiting

| Plan | Requests/minute | Requests/day |
|------|-----------------|--------------|
| Free | 60 | 1,000 |
| Pro | 600 | 50,000 |
| Enterprise | Unlimited | Unlimited |

## Response Format

All responses are JSON:

```json
{
  "data": { ... },
  "meta": {
    "request_id": "req_abc123",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

## Error Handling

Errors follow a consistent format:

```json
{
  "error": {
    "code": "invalid_request",
    "message": "The request was invalid",
    "details": { ... }
  }
}
```

> [!NOTE]
> All timestamps are in UTC and formatted as ISO 8601.
