 # Start from the official Go image
FROM golang:1.21-alpine

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the http port
EXPOSE 8080

# Run the application
CMD ["./main"]