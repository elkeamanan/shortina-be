FROM golang:1.23.4-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migrations directory
COPY --from=builder /app/storage/postgres/migrations ./storage/postgres/migrations/

# Copy any other necessary files (config, templates, etc.)
# COPY --from=builder /app/config ./config/

RUN mkdir -p ./certs
COPY --from=builder /app/ssl-root-cert.pem ./certs/ssl-root-cert.pem

# Copy environment variable
ARG ENV
RUN test -n "$ENV" || (echo "ERROR: ENV build argument is required. Use --build-arg ENV=staging or ENV=release" && exit 1)
COPY .env.${ENV} .env

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Command to run the application
CMD ["./main"]