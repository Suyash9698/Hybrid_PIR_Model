FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/client
RUN go build -o client .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/client/client .
CMD ["./client"]
