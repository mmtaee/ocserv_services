FROM ubuntu:24.04

RUN apt update && apt install -y --no-install-recommends ocserv gnutls-bin build-essential iptables cron openssl sudo && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir /app

COPY deploy/entrypoint.sh /entrypoint.sh

COPY deploy/start.sh /start.sh

RUN chmod +x /entrypoint.sh /start.sh

EXPOSE 443/tcp 443/udp

VOLUME ["/etc/ocserv", "/app"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/start.sh"]

