# Use an official Golang image as the base image
FROM golang:1.20-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the workspace
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app with optimizations for a smaller binary
RUN go build -o main ./cmd/auth-service/main.go

# Use a smaller base image to reduce image size
FROM alpine:latest

# Set up working directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder image
COPY --from=build /app/main ./

# Copy configuration files (like config.yaml) from the source code
COPY --from=build /app/configs /app/configs

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
