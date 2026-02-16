FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/migrate .
COPY --from=builder /app/internal/config/config.yaml ./internal/config/config.yaml
CMD ["./migrate"]
