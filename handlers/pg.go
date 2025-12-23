package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getPrices(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := pool.Query(ctx, `
    SELECT product_id, created_at, name, category, price
    FROM prices
    ORDER BY product_id ASC
  `)
	if err != nil {
		http.Error(w, "db query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var csvBuf bytes.Buffer
	cw := csv.NewWriter(&csvBuf)
	cw.Comma = ','

	for rows.Next() {
		var (
			id        int64
			createdAt time.Time
			name      string
			category  string
			price     float64
		)

		if err := rows.Scan(&id, &createdAt, &name, &category, &price); err != nil {
			http.Error(w, "db scan failed", http.StatusInternalServerError)
			return
		}

		rec := []string{
			strconv.FormatInt(id, 10),
			createdAt.Format("2006-01-02"),
			name,
			category,
			formatPrice2(price),
		}

		if err := cw.Write(rec); err != nil {
			http.Error(w, "csv write failed", http.StatusInternalServerError)
			return
		}
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "db rows error", http.StatusInternalServerError)
		return
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		http.Error(w, "csv flush failed", http.StatusInternalServerError)
		return
	}

	zipBytes, err := buildZipWithDataCSV(csvBuf.Bytes())
	if err != nil {
		http.Error(w, "zip build failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="data.zip"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(zipBytes)
}

func formatPrice2(v float64) string {
	return fmt.Sprintf("%.2f", v)
}
