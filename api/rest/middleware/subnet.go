package middleware

import (
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/subnet"
	"net/http"
)

// CheckSubnet middleware.
func CheckSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sn := config.GetConfig().TrustedSubnet
		ip := r.Header.Get("X-Real-IP")

		if !subnet.CheckSubnet(ip, sn) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
