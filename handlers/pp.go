package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func postPrices(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body []byte

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		// tests.sh шлёт файл именно так: -F "file=@..."
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "cannot parse multipart", http.StatusBadRequest)
			return
		}

		f, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file field is required", http.StatusBadRequest)
			return
		}
		defer f.Close()

		body, err = io.ReadAll(f)
		if err != nil {
			http.Error(w, "cannot read uploaded file", http.StatusBadRequest)
			return
		}
	} else {
		// старый режим: zip прямо в body
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}
	}
	csvRC, err := getCSVFromZipBody(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer csvRC.Close()

	reader := csv.NewReader(csvRC)
	reader.Comma = ','

	first, err := reader.Read()
	if err != nil {
		http.Error(w, "bad csv", http.StatusBadRequest)
		return
	}

	var idx colIndex

	if looksLikeHeader(first) {
		idx = parseHeader(first)
	} else {
		idx = detectOrderByData(first)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		http.Error(w, "db begin failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	if !looksLikeHeader(first) {
		if err := insertRow(ctx, tx, first, idx); err != nil {
			http.Error(w, "db insert failed", http.StatusInternalServerError)
			return
		}
	}

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "bad csv", http.StatusBadRequest)
			return
		}

		if err := insertRow(ctx, tx, rec, idx); err != nil {
			http.Error(w, "db insert failed", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "db commit failed", http.StatusInternalServerError)
		return
	}

	var totalItems int64
	var totalCategories int64
	var totalPrice float64

	err = pool.QueryRow(ctx, `
    SELECT
      COUNT(*),
      COUNT(DISTINCT category),
      COALESCE(SUM(price)::double precision, 0)
    FROM prices
  `).Scan(&totalItems, &totalCategories, &totalPrice)
	if err != nil {
		http.Error(w, "db stats failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"total_items":      totalItems,
		"total_categories": totalCategories,
		"total_price":      totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

type colIndex struct {
	id       int
	date     int
	name     int
	category int
	price    int
}

func looksLikeHeader(rec []string) bool {
	if len(rec) == 0 {
		return false
	}
	first := strings.ToLower(strings.TrimSpace(rec[0]))
	if first == "id" || first == "product_id" {
		return true
	}
	joined := strings.ToLower(strings.Join(rec, ","))
	return strings.Contains(joined, "price") && strings.Contains(joined, "category")
}

func parseHeader(header []string) colIndex {
	idx := colIndex{id: 0, name: 1, category: 2, price: 3, date: 4}

	m := map[string]int{}
	for i, col := range header {
		key := strings.ToLower(strings.TrimSpace(col))
		m[key] = i
	}

	if v, ok := m["id"]; ok {
		idx.id = v
	} else if v, ok := m["product_id"]; ok {
		idx.id = v
	}

	if v, ok := m["name"]; ok {
		idx.name = v
	}

	if v, ok := m["category"]; ok {
		idx.category = v
	}

	if v, ok := m["price"]; ok {
		idx.price = v
	}

	if v, ok := m["create_date"]; ok {
		idx.date = v
	} else if v, ok := m["created_at"]; ok {
		idx.date = v
	} else if v, ok := m["created_date"]; ok {
		idx.date = v
	}

	return idx
}
func detectOrderByData(first []string) colIndex {
	idxA := colIndex{id: 0, name: 1, category: 2, price: 3, date: 4}
	idxB := colIndex{id: 0, date: 1, name: 2, category: 3, price: 4}

	if len(first) >= 5 {
		if _, err := time.Parse("2006-01-02", strings.TrimSpace(first[1])); err == nil {
			return idxB
		}
		if _, err := time.Parse("2006-01-02", strings.TrimSpace(first[4])); err == nil {
			return idxA
		}
	}

	return idxA
}

func insertRow(ctx context.Context, tx pgx.Tx, rec []string, idx colIndex) error {
	idStr := strings.TrimSpace(rec[idx.id])
	name := rec[idx.name]
	category := rec[idx.category]
	priceStr := strings.TrimSpace(rec[idx.price])
	dateStr := strings.TrimSpace(rec[idx.date])

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return err
	}

	createdAt, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO prices (product_id, created_at, name, category, price)
     VALUES ($1, $2, $3, $4, $5)`,
		id, createdAt, name, category, price,
	)
	return err
}
