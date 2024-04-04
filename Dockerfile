# Use the official golang image as the base
FROM golang:latest

# Set a working directory for the application
WORKDIR /app

# Copy the Go source code (including go.mod and go.sum)
COPY . .

# Download Go dependencies
RUN go mod download

# Build the Go binary (assuming main executable is in main.go)
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Set the command to run the application (entrypoint)
CMD ["./main"]