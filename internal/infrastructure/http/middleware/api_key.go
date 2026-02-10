package middleware

import (
	"net/http"
	"os"
)

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("apiKey")

		expectedAPIKey := os.Getenv("API_KEY")

		if apiKey == "" || apiKey != expectedAPIKey {
			http.Error(w, "Acesso negado: API Key inválida", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
