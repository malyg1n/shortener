package middleware

import (
	"github.com/malyg1n/shortener/pkg/config"
	"net"
	"net/http"
)

// CheckSubnet middleware.
func CheckSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subnet := config.GetConfig().TrustedSubnet
		if subnet == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		_, sNet, err := net.ParseCIDR(subnet)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		sIP := net.ParseIP(ip)
		if !sNet.Contains(sIP) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
