FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ingestion ./cmd/ingestion/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ingestion .
CMD ["./ingestion"]
