# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s -X main.Version=$(VERSION) -X main.Commit=$(COMMIT)" \
  -o /bin/vault-agent \
  ./cmd/vault-agent

# Final stage
FROM alpine:latest

# Install necessary tools
RUN apk add --no-cache ca-certificates

# Copy the built binary
COPY --from=builder /bin/vault-agent /bin/vault-agent

# Set user to non-root
USER nobody

# Set entrypoint
ENTRYPOINT ["/bin/vault-agent"]
