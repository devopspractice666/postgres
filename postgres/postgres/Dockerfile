FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o myapp .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp .

EXPOSE 9055

CMD ["./myapp"]