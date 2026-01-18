#!/bin/sh
set -e

echo "Running migrations..."
./migrations

echo "Starting application..."
exec "$@"
