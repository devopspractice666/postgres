FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o myapp .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -g 65532 -S app \
    && adduser -S -D -H -u 65532 -G app app

WORKDIR /app

COPY --from=builder --chown=65532:65532 /app/myapp .

USER app

EXPOSE 9055

HEALTHCHECK --interval=30s --timeout=3s --start-period=15s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:9055/health || exit 1

CMD ["./myapp"]
