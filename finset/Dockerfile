# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /finset \
    ./cmd/server

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.20

# CA certs for TLS connections to PostgreSQL
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /finset /app/finset

# Railway injects PORT automatically
EXPOSE 8080

CMD ["/app/finset"]
