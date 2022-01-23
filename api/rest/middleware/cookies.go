package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/pkg/crypto"
	"net/http"
)

func Cookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_uuid")
		if err != nil {
			userUUID := uuid.New().String()
			encrypted, err := crypto.Encrypt(userUUID)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), "user_uuid", userUUID))
				http.SetCookie(w, &http.Cookie{Name: "user_uuid", Value: encrypted, Path: "/"})
			}
		} else {
			decrypted, err := crypto.Decrypt(cookie.Value)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), "user_uuid", decrypted))
			}
		}

		next.ServeHTTP(w, r)
	})
}
