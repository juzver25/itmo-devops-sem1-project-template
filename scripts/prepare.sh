#!/usr/bin/env bash
set -euo pipefail

sudo apt update
sudo apt install -y curl zip unzip jq

sudo apt install -y postgresql postgresql-contrib

sudo systemctl start postgresql
sudo systemctl enable postgresql >/dev/null 2>&1  true


sudo -u postgres psql <<'SQL'
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'validator') THEN
    CREATE ROLE validator LOGIN PASSWORD 'val1dat0r';
  END IF;
END$$;
SQL

sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='project-sem-1'" | grep -q 1 \
   sudo -u postgres psql -c 'CREATE DATABASE "project-sem-1" OWNER validator;'

sudo -u postgres psql -d "project-sem-1" <<'SQL'
CREATE TABLE IF NOT EXISTS prices (
  product_id  BIGINT,
  created_at  DATE,
  name        TEXT,
  category    TEXT,
  price       DOUBLE PRECISION
);
ALTER TABLE prices OWNER TO validator;
GRANT ALL PRIVILEGES ON TABLE prices TO validator;
SQL

go mod tidy

echo "O"
