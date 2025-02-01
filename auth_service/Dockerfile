# Use the official Golang image as a builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files first (for caching dependencies)
COPY go.mod ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o auth-service main.go

# Use a minimal image to run the binary
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/auth-service .

# Expose the required port
EXPOSE 8081

# Set environment variables (can be overridden at runtime)
ENV PORT=${PORT}
ENV JWT_SECRET=${JWT_SECRET}

# Run the binary
CMD ["./auth-service"]
