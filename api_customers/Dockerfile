FROM golang:1.23 AS builder
LABEL authors="masoud"

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o api_customers main.go

FROM debian:bullseye-slim

COPY --from=builder /app/api_customers /api_customers

RUN chmod +x /api_customers

CMD ["sh", "-c", "/api_customers"]
