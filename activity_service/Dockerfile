# Use the official Golang image for building
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go modules and install dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o activity-service ./main.go

# Use a lightweight image for deployment
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/activity-service .

# Expose the application's port
EXPOSE 8080

# Start the service
CMD ["./activity-service"]
