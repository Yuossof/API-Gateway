package ratelimit

import (
	_ "embed"

	"net/http"

	// "github.com/Yuossof/go-api-gateway/pkg/response"

	"github.com/Yuossof/go-api-gateway/internal/gatewayconfig"
	"github.com/Yuossof/go-api-gateway/pkg/response"
	"github.com/redis/go-redis/v9"
)

const clientIPKey string = "clientIP"

const MaxTimeForToken int64 = 6

func NewRateLimiter(
	client *redis.Client,
	rateLimit gatewayconfig.RateLimit,
) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			ip, ok := r.Context().Value(clientIPKey).(string)

			if !ok {
				http.Error(w, "Client IP not found", http.StatusInternalServerError)
				return
			}

			allowed, err := rateLimitScript.Run(
				ctx,
				client,
				[]string{ip},
				rateLimit.NumberOfTokens,
				rateLimit.MaxTimeForToken,
			).Int()

			if err != nil {
				http.Error(
					w,
					"Rate limiter unavailable",
					http.StatusServiceUnavailable,
				)
				return
			}

			if allowed == 0 {
				response.TooManyRequests(w, "Too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}

//go:embed redis_lua_script.lua
var luaScript string

var rateLimitScript = redis.NewScript(luaScript)
