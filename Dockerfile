# Start with a base image that has Go installed
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o auth-service ./cmd/auth-service/main.go

# Ensure the binary is executable
RUN chmod +x /app/auth-service

# Start a new stage from scratch
FROM alpine:latest AS runtime

# Install necessary dependencies
RUN apk add --no-cache libc6-compat

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/auth-service .

# Copy the config and db folders into the container
COPY --from=builder /app/config ./config
COPY --from=builder /app/db ./db

# Expose port 8080 (or whatever port your service uses)
EXPOSE 8080

# Command to run the executable
CMD ["./auth-service"]
