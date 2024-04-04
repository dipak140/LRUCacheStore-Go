# Use the official Golang image as a base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY *.go ./

# Build the Go app
RUN go build -o /github.com/dipak140/LRUCacheStore-Go

EXPOSE 8080

# Command to run the executable
CMD ["./github.com/dipak140/LRUCacheStore-Go"]
