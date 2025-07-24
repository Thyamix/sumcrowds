FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY . .

WORKDIR /app/counter

RUN go mod download

RUN go build -o go-backend -ldflags="-s -w" ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/counter/go-backend .

RUN chmod +x /app/go-backend

EXPOSE 8080

CMD ["./go-backend"]
