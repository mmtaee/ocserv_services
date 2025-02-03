#!/bin/bash

/usr/sbin/ocserv --debug=9999 --foreground --config=/etc/ocserv/ocserv.conf &
/ocserv_api -migrate && /ocserv_api &

wait -n
exit $?


