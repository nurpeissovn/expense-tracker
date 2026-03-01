# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (required by go mod tidy)
RUN apk --no-cache add git ca-certificates

# Copy go.mod then download deps
COPY go.mod .
RUN go mod download

# Copy source
COPY . .

# Ensure static assets are colocated for go:embed static/*
RUN mkdir -p cmd/server/static && cp -r static/* cmd/server/static/ || true

# Recreate go.sum (to avoid checksum mismatches)
RUN rm -f go.sum && go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /finset ./cmd/server

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /finset /app/finset
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
ENTRYPOINT ["/app/finset"]
