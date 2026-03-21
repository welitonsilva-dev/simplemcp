# stage 1 — build
FROM golang:1.25 AS builder

WORKDIR /app

COPY humancli-server/ ./humancli-server/
COPY humancli-plugins/ ./humancli-plugins/

WORKDIR /app/humancli-server

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o humancli-server ./cmd/server

# stage 2 — runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/humancli-server/humancli-server .

EXPOSE 8081
CMD ["./humancli-server"]