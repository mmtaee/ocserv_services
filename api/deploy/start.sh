#!/bin/bash

/usr/sbin/cron -f &
/usr/sbin/ocserv --debug=9999 --foreground --config=/etc/ocserv/ocserv.conf >> /var/log/ocserv/ocserv.log 2>&1 &
/ocserv_api -migrate && /ocserv_api &

wait -n
exit $?


