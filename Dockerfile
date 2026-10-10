# ==========================================
# Build stage
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary from cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/app

# ==========================================
# Runtime stage
# ==========================================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /out/app /app/app
COPY migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/app"]
