package middleware

import (
	"context"
	"fmt"
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
		fmt.Println("cookies starting")
		cookie, err := r.Cookie(userKey)
		fmt.Println("cookie err", err)
		if err != nil {
			userUUID := uuid.New().String()
			encrypted, err := crypto.Encrypt(userUUID)
			fmt.Println("new cookies starting")
			if err == nil {
				fmt.Println("set context")
				r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, userUUID))
				http.SetCookie(w, &http.Cookie{Name: userKey, Value: encrypted, Path: "/"})
			}
			fmt.Println("encrypt cookie err", err)
		} else {
			decrypted, err := crypto.Decrypt(cookie.Value)
			fmt.Println("exists cookies starting")
			if err == nil {
				fmt.Println("set context")
				r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, decrypted))
			}
			fmt.Println("decrypt cookie err", err)
		}

		next.ServeHTTP(w, r)
	})
}
