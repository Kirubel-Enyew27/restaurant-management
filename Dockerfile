# Use the official Golang image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire cmd directory to the working directory
COPY cmd/ ./cmd/

# Copy any other necessary files (e.g., if you have a pkg or internal directory)
COPY internal/ ./internal/  
COPY intiator/ ./intiator/  
COPY config/ ./config/  
COPY utils/ ./utils/  

# Copy the migrations directory.
COPY internal/constant/query/schemas  /app/internal/constant/query/schemas

# Build the Go application
RUN go build -o myapp ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/myapp .

# Expose port 8081 for the app
EXPOSE 8081

# Command to run the executable
CMD ["./myapp"]
