# Build stage
FROM golang:1.22-alpine AS builder
# Set working directory
WORKDIR /app
# Install git for fetching dependencies
RUN apk add --no-cache git
# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum* ./
# Download dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o discordbot .

# Final stage
FROM alpine:latest
# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/discordbot .
# Copy any config files if needed
COPY --from=builder /app/config* ./
# Copy the .env file
COPY .env ./
# Copy the sounds directory
COPY sounds/ ./sounds/
# Expose port if your bot needs to listen on a port
# EXPOSE 8080
# Command to run the executable
CMD ["./discordbot"]
