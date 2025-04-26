# Use a base image with Go pre-installed
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o pdf-service

# Use a smaller base image for the final build, including Chromium
FROM ubuntu:22.04

# Install Chromium and dependencies
RUN apt-get update && apt-get install -y \
    chromium-browser \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libexpat1 \
    libfontconfig1 \
    libgbm1 \
    libgdk-pixbuf-2.0-0 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    && rm -rf /var/lib/apt/lists/*

# Copy the Go binary from the builder stage
COPY --from=builder /app/pdf-service /app/pdf-service

# Copy static files (fonts and logo)
COPY static /app/static

# Set working directory
WORKDIR /app

# Expose port 8080
EXPOSE 8080

# Command to run the service
CMD ["./pdf-service"]