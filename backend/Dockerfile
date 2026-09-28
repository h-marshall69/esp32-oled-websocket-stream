# Dockerfile

# =====================================================
# BUILD
# =====================================================

FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    git \
    ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/oled-server \
    ./cmd/server


# =====================================================
# RUNTIME
# =====================================================

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

COPY --from=builder /bin/oled-server /app/oled-server

USER app

EXPOSE 8080

CMD ["/app/oled-server"]
