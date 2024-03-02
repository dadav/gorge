package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func AuthMiddleware(tokenAuth *jwtauth.JWTAuth, isProtected func(*http.Request) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isProtected(r) {
				jwtauth.Verifier(tokenAuth)(jwtauth.Authenticator(tokenAuth)(next)).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
