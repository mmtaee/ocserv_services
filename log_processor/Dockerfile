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

RUN go build -o log_processor cmd/main.go

FROM alpine:latest

RUN mkdir -p /app /var/ocserv

COPY --from=builder /app/log_processor /log_processor

RUN chmod +x /log_processor

VOLUME ["/var/log/ocserv"]

CMD ["sh", "-c", "/log_processor -file"]
