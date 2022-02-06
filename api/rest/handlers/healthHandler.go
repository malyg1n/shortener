package handlers

import (
	_ "github.com/lib/pq"
	"github.com/malyg1n/shortener/storage/pgsql"
	"net/http"
)

// PingDB checks connection to DB.
func (hm *HandlerManager) PingDB(w http.ResponseWriter, r *http.Request) {

	storage, err := pgsql.NewLinksStoragePG(r.Context())
	if err != nil {
		http.Error(w, "Db connection refused", http.StatusInternalServerError)
		return
	}

	err = storage.Ping()
	if err != nil {
		http.Error(w, "Db connection refused", http.StatusInternalServerError)
		return
	}

	storage.Close()

	w.WriteHeader(http.StatusOK)
}
