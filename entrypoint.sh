#!/bin/sh
set -e

if [ "$1" = 'alterego' ]; then
    /usr/bin/alterego
fi

exec "$@"%