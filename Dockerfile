# Use the official Golang image as the base image
FROM golang:alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the dependency files to the working directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the source files to the working directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Staging the final image
FROM alpine

COPY --from=builder /app/main /
COPY templates /templates

# ENVIRONMENT VARIABLES
ENV GIN_MODE=release

# Expose port 8080
EXPOSE 8080

# Set the command to run the executable
ENTRYPOINT ["/main"]