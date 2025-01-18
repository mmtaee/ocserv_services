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

RUN go build -o user_expiry main.go

FROM alpine:latest

RUN apk update && apk add --no-cache dcron && rm -rf /var/cache/apk/*

COPY --from=builder /app/user_expiry /user_expiry

RUN mkdir -p /var/log/cron

RUN chmod +x /user_expiry

RUN echo "5 0 * * * root /user_expiry -restore >> /var/log/cron.log 2>&1" >> /etc/crontab && \
    echo "0 6,12,18,23 * * * root /user_expiry -expire >> /var/log/cron.log 2>&1" >> /etc/crontab

CMD ["sh", "-c", "crond -f"]
