# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-service

# Final stage
FROM scratch

WORKDIR /

COPY --from=builder /api-service /api-service
COPY --from=builder /app/public /public

EXPOSE 3001

USER 10001

ENTRYPOINT ["/api-service"]
