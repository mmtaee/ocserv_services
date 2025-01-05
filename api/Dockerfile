FROM ubuntu:24.04

ENV SECRET_KEY_FILE_NAME=/tmp/init_secret

RUN apt update &&\
    apt install -y --no-install-recommends ocserv gnutls-bin build-essential iptables cron openssl sudo logrotate &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir /app

COPY deploy/entrypoint.sh /entrypoint.sh

COPY deploy/start.sh /start.sh

RUN chmod +x /entrypoint.sh /start.sh

RUN SECRET=$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32) && \
    echo "${SECRET}" > ${SECRET_KEY_FILE_NAME}

EXPOSE 443/tcp 443/udp

VOLUME ["/etc/ocserv", "/app", "/var/log/"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/start.sh"]

