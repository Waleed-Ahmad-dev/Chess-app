# Build Stage
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker cache
# We use a wildcard for go.sum in case it doesn't exist yet
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the binary
# CGO_ENABLED=0 ensures a static binary for the scratch/alpine container
# -o chess-app names the binary
# ./cmd/chess/main.go is the entry point
RUN CGO_ENABLED=0 GOOS=linux go build -o chess-app ./cmd/chess/main.go

# Run Stage
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/chess-app .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application in web mode
CMD ["./chess-app", "web"]