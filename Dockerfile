# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the final image
FROM alpine:latest

# Install Chromium and dependencies for chromedp
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    fontconfig

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the templates directory
COPY --from=builder /app/templates /app/templates

# Set the working directory
WORKDIR /app

# Set environment variable for Chromium path
ENV CHROME_PATH=/usr/bin/chromium-browser

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]