# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (required by go mod tidy)
RUN apk --no-cache add git

# Copy go.mod only — then download deps (no go.sum yet so no checksum errors)
COPY go.mod ./
RUN go mod download

# Copy all source files
COPY . .

# Make static assets available adjacent to cmd/server for go:embed static/*
RUN mkdir -p cmd/server/static && cp -r static/* cmd/server/static/ || true

# Regenerate go.sum with correct real checksums (fixes any mismatch from local go.sum)
RUN rm -f go.sum && go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /finset \
    ./cmd/server

# ── Runtime stage ───────────────────────────
FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /finset /app/finset
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/  # ensure TLS works

WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/finset"]

