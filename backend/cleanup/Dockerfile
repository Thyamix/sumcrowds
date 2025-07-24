FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY . .

WORKDIR /app/cleanup

RUN go mod download

RUN go build -o go-cleanup -ldflags="-s -w" ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/cleanup/go-cleanup .

RUN chmod +x /app/go-cleanup

EXPOSE 8080

CMD ["./go-cleanup"]
