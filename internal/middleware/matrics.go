package middleware

import (
	"net/http"

	"github.com/Yuossof/go-api-gateway/internal/matrics"
)

func Matrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matrics.RequestsTotal.WithLabelValues(r.Method).Inc()
		next.ServeHTTP(w, r)
	})
}
