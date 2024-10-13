# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# Copy the source code into the container
COPY . .

User root
# Set permissions to ensure the directory is writable
RUN mkdir -p /app && chmod -R 775 /app

# Build the Go app
RUN go build -o main .

# Check if the binary was created
RUN ls -l ./main

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
