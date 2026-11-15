# syntax=docker/dockerfile:1

FROM golang:1.27.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /server \
    ./cmd/gateway

FROM alpine:3.22

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /server ./server
<<<<<<< HEAD
COPY --from=builder /app/configs/gateway/versions ./configs/gateway/versions
=======
>>>>>>> 6f82c121fd53d7d42d588dfbe6d6c0b7fdc3687f

USER appuser

ENTRYPOINT ["./server"]