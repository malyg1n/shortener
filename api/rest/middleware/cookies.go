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
		var userID string

		cookie, err := r.Cookie(userKey)
		if err != nil {
			userID = uuid.New().String()
			var encrypted string
			encrypted, err = crypto.Encrypt(userID)
			if err != nil {
				http.Error(w, "something went wrong", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: userKey, Value: encrypted, Path: "/"})
		} else {
			userID, err = crypto.Decrypt(cookie.Value)
			if err != nil {
				http.Error(w, "something went wrong", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContextUserKey, userID)))
	})
}
