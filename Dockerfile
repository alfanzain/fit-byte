# First stage: Build the Go binary
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go binary for Linux (target architecture)
RUN GOOS=linux GOARCH=amd64 go build -o fit-byte main.go

# Second stage: Create a smaller final image
FROM debian:bookworm-slim

# Install dependencies for running the Go binary (e.g., PostgreSQL client, etc.)
RUN apt-get update && apt-get install -y libpq-dev

# Set the working directory inside the final container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/fit-byte .

# Copy the .env file
COPY .env .

# Expose the port your app will run on
EXPOSE 8080

# Command to run your Go binary
CMD ["./fit-byte"]
