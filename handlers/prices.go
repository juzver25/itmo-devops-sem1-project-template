package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Prices(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postPrices(pool, w, r)
		case http.MethodGet:
			getPrices(pool, w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
