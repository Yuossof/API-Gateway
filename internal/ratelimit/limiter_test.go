package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yuossof/go-api-gateway/internal/gatewayconfig"
	"github.com/redis/go-redis/v9"
)

func TestRateLimiterWithoutClientIP(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	middleware := NewRateLimiter(client, gatewayconfig.RateLimit{})
	handler := middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called")
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d",
			http.StatusInternalServerError, recorder.Code)
	}
}
