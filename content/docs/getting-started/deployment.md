---
title: "Deployment"
description: "Deploy your site to production"
order: 5
---

This guide covers various deployment options for your site.

## Build for Production

First, build your site:

```bash
./site build --base-url https://your-domain.com
```

This generates static files in the `dist/` directory.

## Deployment Options

### Option 1: Static Hosting

For static-only sites (no reactions):

1. Build your site
2. Upload `dist/` to any static host:
   - Netlify
   - Vercel
   - GitHub Pages
   - Cloudflare Pages

### Option 2: Self-Hosted Server

For full functionality including reactions:

```bash
# Production mode
./site serve --port 8080 --base-url https://your-domain.com
```

> [!NOTE]
> The server handles both static files and the reactions API.

### Option 3: Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o site ./cmd/site

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/site .
COPY content/ content/
COPY templates/ templates/
COPY static/ static/
EXPOSE 8080
CMD ["./site", "serve", "--port", "8080"]
```

Build and run:

```bash
docker build -t my-site .
docker run -p 8080:8080 my-site
```

## Environment Configuration

Set these environment variables in production:

```bash
export GOOGLE_CLIENT_ID=your-client-id
export GOOGLE_CLIENT_SECRET=your-client-secret
```

> [!WARNING]
> Always use HTTPS in production for OAuth to work correctly.

## Health Checks

The server exposes a health endpoint:

```bash
curl http://localhost:8080/api/health
```

## Congratulations!

You've completed the Getting Started guide. Your site should now be up and running!

For more advanced topics, check out the other documentation sections.
