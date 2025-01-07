#!/bin/bash

/usr/sbin/cron -f &
/usr/sbin/ocserv --debug=2 --foreground --config=/etc/ocserv/ocserv.conf >/var/log/ocserv/ocserv.log 2>&1 &
/app/ocserv_api -migrate && /app/ocserv_api &

wait -n
exit $?


