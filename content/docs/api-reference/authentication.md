---
title: "[TEST] Authentication"
description: "How to authenticate with the API"
order: 2
---

Learn how to authenticate your API requests.

## Getting an API Key

1. Log in to your dashboard
2. Navigate to Settings → API Keys
3. Click "Create New Key"
4. Copy your key (it won't be shown again!)

> [!WARNING]
> Keep your API keys secret. Never commit them to version control or share them publicly.

## Using Your Key

Include your API key in the `Authorization` header:

```bash
curl -X GET "https://api.example.com/v1/users" \
     -H "Authorization: Bearer sk_live_abc123..."
```

## Key Types

### Live Keys

- Prefix: `sk_live_`
- Use in production
- Full access to your account

### Test Keys

- Prefix: `sk_test_`
- Use in development
- Isolated test environment

```javascript
// Example: Initializing the client
const client = new ApiClient({
  apiKey: process.env.API_KEY,
  environment: "production", // or 'sandbox'
});
```

## OAuth 2.0

For user-facing applications, we support OAuth 2.0:

```
GET /oauth/authorize
  ?client_id=YOUR_CLIENT_ID
  &redirect_uri=https://yourapp.com/callback
  &response_type=code
  &scope=read write
```

## Token Refresh

Access tokens expire after 1 hour. Use refresh tokens to get new ones:

```bash
curl -X POST "https://api.example.com/v1/oauth/token" \
     -d "grant_type=refresh_token" \
     -d "refresh_token=YOUR_REFRESH_TOKEN"
```
