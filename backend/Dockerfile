# Start from official Golang image
FROM golang:1.24-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy everything else
COPY . .

# Copy .env file
COPY .env .env

# Build the Go binary
RUN go build -o main .

# Expose port
EXPOSE 8080

# Run the Go binary
CMD ["./main"]
