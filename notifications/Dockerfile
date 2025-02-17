# Stage 1: Build
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
ARG VERSION
RUN echo "$VERSION" > /app/version.txt
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app

# Stage 2: Run
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy built binary and version file from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/version.txt .

# Expose necessary ports
EXPOSE 8085

# Run the application
CMD ["./main"]
