# stage 1 — build
FROM golang:1.25 AS builder

WORKDIR /app

# Copia arquivos de módulo primeiro para melhor cache
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mcp-server ./cmd

# stage 2 — runtime
FROM debian:bookworm-slim

# Instala certificados necessários para conexões externas
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/mcp-server .

EXPOSE 8081
CMD ["./mcp-server"]