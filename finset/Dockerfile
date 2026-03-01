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
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /finset /app/finset

EXPOSE 8080
ENTRYPOINT ["/app/finset"]
