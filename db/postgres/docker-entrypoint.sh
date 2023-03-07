#!/bin/sh

/usr/bin/sudo /usr/sbin/sshd -D&

exec "$@"