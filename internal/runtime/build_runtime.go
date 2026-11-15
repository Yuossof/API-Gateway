package runtime

import (
	"context"
	"net/http"

	"github.com/Yuossof/go-api-gateway/internal/gatewayconfig"
	"github.com/Yuossof/go-api-gateway/internal/health"
	"github.com/Yuossof/go-api-gateway/internal/loadbalancer"
	"github.com/Yuossof/go-api-gateway/internal/middleware"
	"github.com/Yuossof/go-api-gateway/internal/proxy"
	"github.com/Yuossof/go-api-gateway/internal/ratelimit"
	"github.com/Yuossof/go-api-gateway/internal/router"
	"github.com/Yuossof/go-api-gateway/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Runtime struct {
	handler http.Handler
	cancel  context.CancelFunc
}

func BuildRuntime(gatewayConfig *gatewayconfig.Config, redis *redis.Client) (*Runtime, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Load proxies
	proxies, err := proxy.LoadProxies(gatewayConfig.Clusters)
	if err != nil {
		// slog.Error("Failed to load proxies", "error", err)
		// os.Exit(1)

		cancel()
		return nil, err
	}

	// Create Health Checker
	healthChecker := health.NewHealthChecker(gatewayConfig.AllServersByCluster)
	healthChecker.CheckOnce(ctx)
	go healthChecker.CheckHealth(ctx)

	// Create Loadbalancer
	lb := loadbalancer.NewLoadBalancer()

	// Create Ratelimiter
	rateLimiterHandler := ratelimit.NewRateLimiter(redis, gatewayConfig.RateLimit)

	routerHandler := router.Router(proxies, gatewayConfig, healthChecker, lb)

	logger := middleware.Logging(logger.Log)

	finalHandler := logger(
		middleware.Matrics(
			middleware.Recovery(
				middleware.IPM(
					rateLimiterHandler(routerHandler),
				),
			),
		),
	)

	return &Runtime{
		handler: finalHandler,
		cancel:  cancel,
	}, nil
}

func (rt *Runtime) Stop() {
	if rt == nil || rt.cancel == nil {
		return
	}

	rt.cancel()
}
