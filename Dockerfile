# Build Stage
# Using latest alpine to support newer Go versions defined in go.mod
FROM golang:alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS (if needed) and migration tool if we want to run migrations inside container
# For now, just the app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
