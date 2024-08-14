# Use Alpine as the base image for building the Go application
FROM golang:1.22.5-alpine as builder

# Install necessary build dependencies
RUN apk add --no-cache \
    build-base \
    git \
    pkgconfig \
    vips-dev \
    ffmpeg

# Set the working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /app/bin/worker ./cmd/main.go

# Use a lightweight Alpine image for the runtime environment
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    vips \
    ffmpeg

# Set the working directory
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/bin/worker /app/worker

# Command to run the application
CMD ["/app/worker"]
