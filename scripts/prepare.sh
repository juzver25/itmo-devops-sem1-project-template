#!/usr/bin/env bash
set -euo pipefail

sudo apt-get update
sudo apt-get install -y postgresql-client curl zip unzip jq

PGHOST="${POSTGRES_HOST:-localhost}"
PGPORT="${POSTGRES_PORT:-5432}"
PGDATABASE="${POSTGRES_DB:-project-sem-1}"
PGUSER="${POSTGRES_USER:-validator}"
PGPASSWORD="${POSTGRES_PASSWORD:-val1dat0r}"
export PGPASSWORD

for i in {1..40}; do
  if psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -c '\q' >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done

psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" <<'SQL'
CREATE TABLE IF NOT EXISTS prices (
  id BIGINT,
  name TEXT,
  category TEXT,
  price DOUBLE PRECISION,
  create_date DATE
);
SQL

go mod tidy
