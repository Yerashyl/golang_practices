package middleware

import (
	"math/rand"
	"net/http"
	"strconv"
)

func RequestID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := strconv.Itoa(rand.Int())

		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r)
	})
}
