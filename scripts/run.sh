#!/usr/bin/env bash
set -euo pipefail

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
  echo "Cannot detect entry point to run."
  echo "Repo files:"
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
done

echo "Server failed to start in time. Log:"
exit 1
