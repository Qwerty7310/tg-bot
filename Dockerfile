FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bot main.go

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/bot .
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
CMD ["./bot"]
