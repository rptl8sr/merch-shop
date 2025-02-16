# Build stage
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o server ./cmd/api/main.go

# Final stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]