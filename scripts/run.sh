#!/usr/bin/env bash
set -euo pipefail

export POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
export POSTGRES_PORT="${POSTGRES_PORT:-5432}"
export POSTGRES_DB="${POSTGRES_DB:-project-sem-1}"
export POSTGRES_USER="${POSTGRES_USER:-validator}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-val1dat0r}"

export DB_HOST="${DB_HOST:-$POSTGRES_HOST}"
export DB_PORT="${DB_PORT:-$POSTGRES_PORT}"
export DB_NAME="${DB_NAME:-$POSTGRES_DB}"
export DB_USER="${DB_USER:-$POSTGRES_USER}"
export DB_PASSWORD="${DB_PASSWORD:-$POSTGRES_PASSWORD}"

export PORT="${PORT:-8080}"

ENTRY=""

if [ -d "./server" ] && go list -f '{{.Name}}' ./server 2>/dev/null | grep -q '^main$'; then
  ENTRY="./server"
elif [ -d "./cmd/server" ] && go list -f '{{.Name}}' ./cmd/server 2>/dev/null | grep -q '^main$'; then
  ENTRY="./cmd/server"
elif [ -f "./main.go" ]; then
  ENTRY="./main.go"
elif go list -f '{{.Name}}' . 2>/dev/null | grep -q '^main$'; then
  ENTRY="."
else
  ls -la
  echo "Go packages:"
  go list ./...  true
  exit 1
fi

echo "Running: go run $ENTRY"
nohup go run "$ENTRY" > /tmp/app.log 2>&1 &
echo $! > /tmp/app.pid

for i in {1..120}; do
  if curl -fsS "http://localhost:${PORT}/health" >/dev/null 2>&1; then
    echo "Server is up."
    exit 0
  fi

  if ! kill -0 "$(cat /tmp/app.pid)" >/dev/null 2>&1; then
    echo "Server process exited early. Log:"
    tail -n 200 /tmp/app.log  true
    exit 1
  fi

  sleep 0.5
done

echo "Server failed to start in time. Log:"
tail -n 200 /tmp/app.log || true
exit 1
