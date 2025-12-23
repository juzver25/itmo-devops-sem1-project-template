#!/usr/bin/env bash
set -euo pipefail

export POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
export POSTGRES_PORT="${POSTGRES_PORT:-5432}"
export POSTGRES_DB="${POSTGRES_DB:-project-sem-1}"
export POSTGRES_USER="${POSTGRES_USER:-validator}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-val1dat0r}"
export PORT="${PORT:-8080}"

ENTRY="."
if [ -d "./cmd/server" ]; then
  ENTRY="./cmd/server"
elif [ -d "./server" ]; then
  ENTRY="./server"
fi

nohup go run "$ENTRY" > /tmp/app.log 2>&1 &
echo $! > /tmp/app.pid

for i in {1..40}; do
  if curl -fsS "http://localhost:${PORT}/health" >/dev/null 2>&1; then
    exit 0
  fi
  sleep 0.25
done

echo "Server failed to start, log:"
exit 1
