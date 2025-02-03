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

RUN go build -o ocserv_api cmd/main.go

FROM ubuntu:24.04

ENV SECRET_KEY_FILE_NAME=/tmp/init_secret

RUN apt update && apt install -y --no-install-recommends ocserv gnutls-bin build-essential iptables openssl less &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir /app

COPY deploy/entrypoint.sh /entrypoint.sh

COPY deploy/start.sh /start.sh

COPY --from=builder /app/ocserv_api /ocserv_api

RUN chmod +x /entrypoint.sh /start.sh /ocserv_api

RUN SECRET=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32) && \
    echo "${SECRET}" > ${SECRET_KEY_FILE_NAME}

EXPOSE 443/tcp 443/udp

VOLUME ["/etc/ocserv", "/app", "/var/log/ocserv"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/start.sh"]

