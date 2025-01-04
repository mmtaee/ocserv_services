#!/bin/bash

printf "\e[33m########### ocserv service starting ... ###########\e[0m"
printf "\n"
pidfile=/run/ocserv.pid
/usr/sbin/ocserv --debug=2 --foreground --pid-file=${pidfile} --config=/etc/ocserv/ocserv.conf &

printf "\e[33m########### api service starting ... ###########\e[0m"
printf "\n"

/app/ocserv_api -migrate && /app/ocserv_api &

wait -n
exit $?


