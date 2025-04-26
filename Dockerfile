
FROM golang:1.24-alpine AS builder


RUN apk add --no-cache git


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o main .


FROM alpine:latest


RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    fontconfig


COPY --from=builder /app/main /app/main


COPY --from=builder /app/templates /app/templates


WORKDIR /app


ENV CHROME_PATH=/usr/bin/chromium-browser


EXPOSE 8080


CMD ["./main"]