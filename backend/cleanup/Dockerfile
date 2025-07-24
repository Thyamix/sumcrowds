FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o cleanup

FROM alpine:3.19
RUN apk add --no-cache cron
WORKDIR /app
COPY --from=builder /build/cleanup /app/cleanup

COPY cleanup-crontab /etc/crontabs/root

COPY start-cron.sh /app/start-cron.sh
RUN chmod +x /app/start-cron.sh

ENTRYPOINT ["/app/start-cron.sh"]
