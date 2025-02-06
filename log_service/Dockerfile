FROM golang:1.23 AS builder

WORKDIR /app

ENV GIN_MODE=release
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o log_service cmd/main.go

FROM debian:bullseye-slim

ENV DOCKERIZED=true

RUN apt-get update && apt-get install -y docker.io

COPY --from=builder /app/log_service /log_service

RUN chmod +x /log_service

CMD ["sh", "-c", "/log_service"]
