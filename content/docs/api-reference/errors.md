---
title: "Error Handling"
description: "Understanding and handling API errors"
order: 4
---

Learn how to handle errors from the API.

## Error Format

All errors follow this structure:

```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable message",
    "details": {},
    "request_id": "req_abc123"
  }
}
```

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| `200` | Success |
| `201` | Created |
| `400` | Bad Request - Invalid parameters |
| `401` | Unauthorized - Invalid or missing API key |
| `403` | Forbidden - Insufficient permissions |
| `404` | Not Found - Resource doesn't exist |
| `429` | Too Many Requests - Rate limited |
| `500` | Internal Server Error |

## Common Error Codes

### `invalid_request`

The request was malformed or missing required fields.

```json
{
  "error": {
    "code": "invalid_request",
    "message": "Missing required field: email",
    "details": {
      "field": "email",
      "reason": "required"
    }
  }
}
```

### `authentication_failed`

The API key is invalid or expired.

```json
{
  "error": {
    "code": "authentication_failed",
    "message": "Invalid API key"
  }
}
```

### `rate_limit_exceeded`

You've exceeded your rate limit.

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded. Try again in 60 seconds.",
    "details": {
      "retry_after": 60
    }
  }
}
```

> [!TIP]
> Check the `Retry-After` header for when you can make requests again.

## Best Practices

### Retry Logic

Implement exponential backoff for transient errors:

```javascript
async function fetchWithRetry(url, options, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(url, options);
      if (response.status === 429) {
        const retryAfter = response.headers.get('Retry-After') || 60;
        await sleep(retryAfter * 1000);
        continue;
      }
      return response;
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await sleep(Math.pow(2, i) * 1000);
    }
  }
}
```

### Logging

Always log the `request_id` for debugging:

```javascript
try {
  const result = await api.createUser(data);
} catch (error) {
  console.error('API Error:', {
    code: error.code,
    message: error.message,
    requestId: error.request_id
  });
}
```
