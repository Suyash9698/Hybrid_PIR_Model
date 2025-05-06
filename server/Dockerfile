FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/server
RUN go build -o server .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server/server .
# volume for data
VOLUME ["/data"]
EXPOSE 8000
CMD ["./server"]
