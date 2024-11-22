# Stage 1: Build the Go application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project source code
COPY . .

# Build the Go app
RUN go build -o auth-service ./cmd/auth-service/main.go

# Stage 2: Create a minimal runtime environment
FROM alpine:latest AS runtime

# Install necessary runtime dependencies
RUN apk add --no-cache libc6-compat

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go binary from the builder stage
COPY --from=builder /app/auth-service .

# Copy the configs and db directories from the builder stage
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/db ./db

# Expose the port used by the service
EXPOSE 8080

# Command to run the application
CMD ["./auth-service"]
