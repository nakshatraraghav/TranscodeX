FROM golang:1.20-bullseye as builder

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    pkg-config \
    libvips-dev \
    ffmpeg \
    && apt-get clean

# Set the working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /app/worker ./cmd/main.go

FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libvips \
    ffmpeg \
    && apt-get clean

# Set the working directory
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/worker /app/worker

# Command to run the application
CMD ["/app/worker"]
