#!/bin/sh
set -e

if ! whoami &> /dev/null; then
  if [ -w /etc/passwd ]; then
    echo "${USER_NAME:-user}:x:$(id -u):0:${USER_NAME:-user} user:${HOME}:/sbin/nologin" >> /etc/passwd
  fi
fi

exec "$@"
