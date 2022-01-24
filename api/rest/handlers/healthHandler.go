package handlers

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/malyg1n/shortener/pkg/config"
	"net/http"
)

// PingDB checks connection to DB.
func (hm *HandlerManager) PingDB(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		http.Error(w, "Db connection refused", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	err = db.PingContext(r.Context())
	if err != nil {
		http.Error(w, "Db connection refused", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
