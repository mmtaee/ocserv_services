#!/bin/bash

printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
pidfile=
/usr/sbin/ocserv --debug=2 --foreground --config=/etc/ocserv/ocserv.conf >/var/log/ocserv.log 2>&1 &

printf "\e[33m########### api service starting ... ###########\e[0m"
printf "\n"
/app/ocserv_api -migrate && /app/ocserv_api &

wait -n
exit $?


