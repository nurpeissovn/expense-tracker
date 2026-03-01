# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies (context root is repository root)
COPY finset/go.mod finset/go.sum ./
RUN go mod download

# Copy source
COPY finset/. .

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
