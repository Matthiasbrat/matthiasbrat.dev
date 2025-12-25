---
title: "Docker Best Practices I Learned the Hard Way"
description: "Lessons from running containers in production"
date: 2025-01-28
---

After years of running Docker in production, here are the lessons that stuck.

## Use Multi-Stage Builds

Keep images small:

```dockerfile
# Build stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o main .

# Runtime stage
FROM alpine:3.19
COPY --from=builder /app/main /main
CMD ["/main"]
```

> [!TIP]
> This Go image went from 800MB to 15MB using multi-stage builds.

## Don't Run as Root

```dockerfile
FROM node:20-alpine
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --chown=app:app . .
```

## Use .dockerignore

```.dockerignore
node_modules
.git
*.md
.env
```

> [!WARNING]
> Without `.dockerignore`, you might accidentally include secrets or massive directories.

## Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD curl -f http://localhost:8080/health || exit 1
```

## Pin Your Versions

```dockerfile
# Bad
FROM node:latest

# Good
FROM node:20.11.0-alpine3.19
```

> [!IMPORTANT]
> `latest` tags can change unexpectedly and break your builds.
