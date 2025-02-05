FROM golang:1.23 AS builder
LABEL authors="masoud"

ENV GIN_MODE=release
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o process cmd/main.go

FROM debian:bullseye-slim

COPY --from=builder /app/process /process

RUN chmod +x /process

CMD ["sh", "-c", "/process"]
