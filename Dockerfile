# stage 1 — build
FROM golang:1.25 AS builder

WORKDIR /app

COPY simplemcp/ ./simplemcp/
COPY simplemcpplugins/ ./simplemcpplugins/

WORKDIR /app/simplemcp

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mcp-server ./cmd/server

# stage 2 — runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/simplemcp/mcp-server .

EXPOSE 8081
CMD ["./mcp-server"]