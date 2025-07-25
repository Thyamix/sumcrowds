#!/bin/sh

mkdir -p /var/log

touch /var/log/cleanup.log

crond -l 2 -f &

tail -F /var/log/cleanup.log
