package main

import (
  "context"
  "log"
  "net/http"
  "os"
  "time"

  "project_sem/bd"
  "project_sem/handlers"
)

func main() {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  pool, err := db.NewPool(ctx)
  if err != nil {
    log.Fatalf("db connect: %v", err)
  }
  defer pool.Close()

  mux := http.NewServeMux()
  mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
  })

  mux.HandleFunc("/api/v0/prices", handlers.Prices(pool)) 

  addr := ":8080"
  if v := os.Getenv("PORT"); v != "" {
    addr = ":" + v
  }

  log.Printf("listening on %s", addr)
  log.Fatal(http.ListenAndServe(addr, mux))
}
