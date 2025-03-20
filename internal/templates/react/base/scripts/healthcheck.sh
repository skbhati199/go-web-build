#!/bin/sh
set -e

HOST="localhost"
PORT="80"
HEALTH_ENDPOINT="/health"

response=$(curl -s -o /dev/null -w "%{http_code}" http://$HOST:$PORT$HEALTH_ENDPOINT)

if [ "$response" = "200" ]; then
    exit 0
else
    exit 1
fi