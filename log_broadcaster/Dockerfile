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

RUN go build -o log_broadcaster cmd/main.go

FROM alpine:latest

RUN mkdir -p /app /var/ocserv

COPY --from=builder /app/log_broadcaster /log_broadcaster

RUN chmod +x /log_broadcaster

VOLUME ["/var/log/ocserv"]

CMD ["sh", "-c", "/log_broadcaster -file"]