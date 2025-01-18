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

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends cron && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/user_expiry /user_expiry

COPY ./deploy/start.sh /start.sh

RUN mkdir -p /var/log/cron

RUN chmod +x /user_expiry /start.sh

#RUN echo "5 0 * * * root /user_expiry -restore >> /var/log/cron.log 2>&1" >> /etc/crontab && \
#    echo "0 6,12,18,23 * * * root /user_expiry -expire >> /var/log/cron.log 2>&1" >> /etc/crontab

RUN crontab -l | echo "5 0 * * * root /user_expiry -restore >> /var/log/cron.log 2>&1" | crontab - &&\
    crontab -l | echo "0 6,12,18,23 * * * root /user_expiry -expire >> /var/log/cron.log 2>&1" | crontab -

CMD ["/start.sh"]
