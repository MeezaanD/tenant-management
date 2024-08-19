# Build stage
FROM golang:1.20-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["go", "run", "main.go"]
