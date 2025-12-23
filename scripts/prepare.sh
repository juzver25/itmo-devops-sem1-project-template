#!/usr/bin/env bash
set -euo pipefail

sudo systemctl start postgresql
sudo -u postgres psql <<'SQL'
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'validator') THEN
    CREATE ROLE validator LOGIN PASSWORD 'val1dat0r';
  END IF;
END$$;
SQL
sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='project-sem-1'" | grep -q 1 \
|| sudo -u postgres psql -c 'CREATE DATABASE "project-sem-1" OWNER validator;'
sudo -u postgres psql -d "project-sem-1" <<'SQL'
CREATE TABLE IF NOT EXISTS prices (
  product_id  BIGINT,
  created_at  DATE,
  name        TEXT,
  category    TEXT,
  price       BIGINT
);
SQL
echo "OK"
