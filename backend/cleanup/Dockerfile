FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY . .

WORKDIR /app/cleanup

RUN go mod download

RUN go build -o cleanup -ldflags="-s -w" ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/cleanup/cleanup /app/cleanup

COPY --from=builder /app/cleanup/cleanup-crontab /etc/crontabs/root

COPY --from=builder /app/cleanup/start-cron.sh /app/start-cron.sh

RUN chmod +x /app/cleanup

RUN chmod +x /app/start-cron.sh

ENTRYPOINT ["/app/start-cron.sh"]
