# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Production stage
FROM alpine:3
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app", "serve"]
