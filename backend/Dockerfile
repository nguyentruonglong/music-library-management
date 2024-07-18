# Use an official Golang image as the base image
FROM golang:1.21.1

# Install the necessary build tools
RUN apt-get update && apt-get install -y gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache the Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o backend main.go

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./backend"]
