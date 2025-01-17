FROM golang:1.23 AS builder
LABEL authors="masoud"

ENV GIN_MODE=release
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV LOG_FILE_PATH=/var/log/ocserv/ocserv.log

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o service_log cmd/main.go

FROM alpine:latest

RUN mkdir -p /app /var/ocserv

COPY --from=builder /app/service_log /service_log

RUN chmod +x /service_log

VOLUME ["/var/log/ocserv"]

CMD ["sh", "-c", "/service_log -file"]
