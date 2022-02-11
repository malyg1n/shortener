package handlers

import "net/http"

// PingDB checks connection to DB.
func (hm *HandlerManager) PingDB(w http.ResponseWriter, r *http.Request) {

	err := hm.service.PingStorage()
	if err != nil {
		http.Error(w, "Db connection refused", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
