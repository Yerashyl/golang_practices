package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		log.Printf("%s %s %s {{request processed}}",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
		)
	})
}
