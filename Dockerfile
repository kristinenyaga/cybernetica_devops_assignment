FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o server main.go

FROM alpine:3.20

RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8000

ENTRYPOINT ["./server"]