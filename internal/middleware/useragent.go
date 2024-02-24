package middleware

import (
	"encoding/json"
	"net/http"
)

type UserAgentNotSetResponse struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func RequireUserAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			errorResponse := UserAgentNotSetResponse{
				Message: "User-Agent header is missing",
				Errors:  []string{"User-Agent must have some value"},
			}
			jsonError, err := json.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
