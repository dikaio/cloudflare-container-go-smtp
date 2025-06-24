# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o smtp-server main.go

# Final stage
FROM scratch

# Copy the binary from builder
COPY --from=builder /app/smtp-server /smtp-server

# Copy CA certificates for SMTP TLS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose port 8080 (Cloudflare default)
EXPOSE 8080

# Run the application
ENTRYPOINT ["/smtp-server"]