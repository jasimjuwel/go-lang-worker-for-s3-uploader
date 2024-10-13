# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.27.3

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# Copy the source code into the container
COPY . .


# Set permissions to ensure the directory is writable
RUN chmod -R 775 /app

# Build the Go app
#RUN go build -o main .

# Check if the binary was created
#RUN ls -l ./main

# Expose port 8080 to the outside world
#EXPOSE 8080

# Command to run the executable
CMD ["air"]
