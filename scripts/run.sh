#!/usr/bin/env bash
set -euo pipefail

export DB_HOST="${DB_HOST:-localhost}"
export DB_PORT="${DB_PORT:-5432}"
export DB_USER="${DB_USER:-validator}"
export DB_PASSWORD="${DB_PASSWORD:-val1dat0r}"
export DB_NAME="${DB_NAME:-project-sem-1}"
export PORT="${PORT:-8080}"

go run /home/vm1/finalboss/itmo-devops-sem1-project-template/server

