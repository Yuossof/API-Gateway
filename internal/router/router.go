package router

import (
	"net/http"
	"strings"

	"github.com/Yuossof/go-api-gateway/internal/cluster"
	"github.com/Yuossof/go-api-gateway/internal/gatewayconfig"
	"github.com/Yuossof/go-api-gateway/internal/health"
	"github.com/Yuossof/go-api-gateway/internal/loadbalancer"
	"github.com/Yuossof/go-api-gateway/internal/proxy"
	"github.com/Yuossof/go-api-gateway/pkg/response"
)

func MatchRoute(path string, routes []gatewayconfig.Route) (gatewayconfig.Route, bool) {
	for _, route := range routes {
		if path == route.Path || strings.HasPrefix(path, route.Path+"/") {
			return route, true
		}
	}

	return gatewayconfig.Route{}, false
}

func rewritePath(path, routePath, target string) string {
	suffix := strings.TrimPrefix(path, routePath)
	return target + suffix
}

func Router(proxies proxy.Proxies, config *gatewayconfig.Config, healthChecker *health.HealthChecker, loadbalancer *loadbalancer.LoadBalancer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, ok := MatchRoute(r.URL.Path, config.Routes)

		if !ok {
			response.NotFound(w, "Path not found")
			return
		}

		clusters, ok := config.AllServersByCluster[route.Cluster]
		if !ok {
			response.NotFound(w, "Cluster not found")
			return
		}

		healthyServers := healthChecker.GetHealthyServers(clusters)

		if len(healthyServers) == 0 {
			response.JSONError(w, "503 Service Unavailable", 503)
			return
		}

		serv := loadbalancer.NextServer(route.Cluster, serverURLs(healthyServers))

		upstream, ok := proxies[serv]

		if !ok {
			response.NotFound(w, "Server not found")
			return
		}

		r.URL.Path = rewritePath(
			r.URL.Path,
			route.Path,
			route.Target,
		)

		upstream.ServeHTTP(w, r)
	})
}

func serverURLs(servers []cluster.Server) []string {
	urls := make([]string, 0, len(servers))
	for _, server := range servers {
		urls = append(urls, server.URL)
	}
	return urls
}
