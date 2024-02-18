# Use the official Golang image as a base
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /go/src/app

# Copy the Go modules files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/server

# Expose port 8080 to the outside world
EXPOSE 9000

# Run the Go application with environment variables fetched from .env file
CMD ["./main"]
