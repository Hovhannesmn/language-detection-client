# Build stage
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (skip vendor for local replaces)
RUN go mod download -x || true
RUN go mod tidy || true

# Copy source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o language-detection-client ./cmd/client

# Build the mock service (for testing only)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o mock-service ./cmd/mock

# Mock service stage (for testing only)
FROM alpine:latest AS mock

# Install ca-certificates and wget for health checks
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy mock service binary from builder stage
COPY --from=builder /app/mock-service .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port for mock service
EXPOSE 8080

# Health check for mock service
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run the mock service
CMD ["./mock-service"]

# Final stage (main application - no HTTP server)
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/language-detection-client .

# Copy sample test files
COPY --from=builder /app/test_captions.srt* ./

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ./language-detection-client -h >/dev/null 2>&1 || exit 1

# Default command
ENTRYPOINT ["./language-detection-client"]

# Default arguments (can be overridden)
CMD ["-h"]
