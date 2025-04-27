# Stage 1: Build and Test
FROM golang:1.24-alpine AS builder

# Install git and Chromium for tests
RUN apk add --no-cache \
    git \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    fontconfig

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Set Chrome path for chromedp
ENV CHROME_PATH=/usr/bin/chromium-browser

# Run all tests, including integration tests
RUN RUN_INTEGRATION_TESTS=true go test ./...

# Build the binary
RUN go build -o main .

# Stage 2: Final Image
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

# Copy templates (if needed by your application)
COPY --from=builder /app/templates /app/templates

# Set working directory
WORKDIR /app

# Set environment variable for Chrome path
ENV CHROME_PATH=/usr/bin/chromium-browser

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]