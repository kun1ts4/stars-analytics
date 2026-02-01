FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o processor ./cmd/processor/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/processor .
CMD ["./processor"]
