package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/pkg/crypto"
	"net/http"
)

type ContextKey string

const (
	userKey        = "user_uuid"
	ContextUserKey = ContextKey(userKey)
)

func Cookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(userKey)
		if err != nil {
			userUUID := uuid.New().String()
			encrypted, err := crypto.Encrypt(userUUID)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, userUUID))
				http.SetCookie(w, &http.Cookie{Name: userKey, Value: encrypted, Path: "/"})
			}
		} else {
			decrypted, err := crypto.Decrypt(cookie.Value)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, decrypted))
			}
		}

		next.ServeHTTP(w, r)
	})
}
