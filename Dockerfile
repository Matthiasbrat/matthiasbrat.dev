# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o site ./cmd/site

# Build the static site
RUN ./site build

# Final stage - distroless for minimal size with CA certs
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy the binary
COPY --from=builder /app/site /app/site

# Copy built static site
COPY --from=builder /app/dist /app/dist

# Copy site.yml if it exists (optional config)
COPY --from=builder /app/site.y[m]l /app/

# Create data directory for SQLite (will be mounted as volume)
# Note: distroless doesn't have mkdir, so we copy an empty dir
COPY --from=builder /app/data* /app/data/

# Expose port
EXPOSE 3000

# Volume for persistent data (SQLite database)
VOLUME ["/app/data"]

# Run the server
ENTRYPOINT ["/app/site"]
CMD ["serve", "--port", "3000"]
