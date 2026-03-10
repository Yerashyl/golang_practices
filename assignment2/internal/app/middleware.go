package app

import (
	"log"
	"net/http"
	"time"
)

func (a *App) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// Логируем: метод, путь и время выполнения [cite: 272-273]
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey != "my-secret-key" { // Замени на свой ключ [cite: 275]
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
