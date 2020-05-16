#!/bin/bash
set -e

if ! whoami &> /dev/null; then
  if [ -w /etc/passwd ]; then
    echo "${USER_NAME:-go}:x:$(id -u):0:${USER_NAME:-go} user:${HOME}:/sbin/nologin" >> /etc/passwd
  fi
fi

exec "$@"

