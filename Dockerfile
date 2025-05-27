# Stage 1: Build the Go binary
FROM golang:1.24.3-alpine AS build

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Copy Templates
COPY templates/ templates/

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary with executable permissions
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web-analyzer ./cmd

# Set execute permission
RUN chmod +x web-analyzer

# Stage 2: Create a minimal image to run the binary
FROM alpine:latest

# Install CA certificates to enable HTTPS
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the previous stage
COPY --from=build /app/web-analyzer /usr/local/bin/web-analyzer

# Copy config
COPY --from=build /app/config /root/config

# Copy templates
COPY --from=build /app/templates /root/templates

# Expose the port on which the service will run
EXPOSE 8080

# Specify the entry point command to run the binary
CMD ["/usr/local/bin/web-analyzer"]