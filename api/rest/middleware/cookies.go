package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/pkg/crypto"
	"net/http"
)

const userUuid = "user_uuid"

func Cookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(userUuid)
		if err != nil {
			userUUID := uuid.New().String()
			encrypted, err := crypto.Encrypt(userUUID)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), userUuid, userUUID))
				http.SetCookie(w, &http.Cookie{Name: "user_uuid", Value: encrypted, Path: "/"})
			}
		} else {
			decrypted, err := crypto.Decrypt(cookie.Value)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), userUuid, decrypted))
			}
		}

		next.ServeHTTP(w, r)
	})
}
