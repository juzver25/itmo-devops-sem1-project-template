package db

import (
  "context"
  "fmt"
  "os"

  "github.com/jackc/pgx/v5/pgxpool"
)

func getenv(key, def string) string {
  if v := os.Getenv(key); v != "" {
    return v
  }
  return def
}

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
  host := getenv("DB_HOST", "localhost")
  port := getenv("DB_PORT", "5432")
  user := getenv("DB_USER", "validator")
  pass := getenv("DB_PASSWORD", "val1dat0r")
  name := getenv("DB_NAME", "project-sem-1")

  dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name)
  cfg, err := pgxpool.ParseConfig(dsn)
  if err != nil {
    return nil, err
  }

  pool, err := pgxpool.NewWithConfig(ctx, cfg)
  if err != nil {
    return nil, err
  }

  // проверяем, что БД реально доступна
  if err := pool.Ping(ctx); err != nil {
    pool.Close()
    return nil, err
  }

  return pool, nil
}
