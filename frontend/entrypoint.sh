#!/bin/sh

envsubst < /app/src/config.template.js > /app/src/config.js

exec "$@"