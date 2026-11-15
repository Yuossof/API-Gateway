package health

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Yuossof/go-api-gateway/internal/cluster"
)

type HealthChecker struct {
	mu                  sync.RWMutex
	status              map[string]bool
	allServersByCluster map[string][]cluster.Server
	client              *http.Client
}

func NewHealthChecker(allServersByCluster map[string][]cluster.Server) *HealthChecker {
	return &HealthChecker{
		status:              map[string]bool{},
		allServersByCluster: allServersByCluster,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (h *HealthChecker) GetHealthyServers(servers []cluster.Server) []cluster.Server {
	healthyServers := make([]cluster.Server, 0, len(servers))

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, s := range servers {
		healthy, ok := h.status[s.URL]
		if !ok {
			continue
		}

		if healthy {
			healthyServers = append(healthyServers, s)
		}

	}

	return healthyServers
}

func (h *HealthChecker) SendReq(ctx context.Context, url string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (h *HealthChecker) checkAll(ctx context.Context) {
	var wg sync.WaitGroup
	for _, servers := range h.allServersByCluster {
		for _, server := range servers {
			wg.Add(1)
			go func(server cluster.Server) {
				defer wg.Done()
				healthURL := strings.TrimRight(server.URL, "/") + "/health"
				healthy := h.SendReq(ctx, healthURL)
				h.mu.Lock()
				defer h.mu.Unlock()
				h.status[server.URL] = healthy
			}(server)
		}
	}

	wg.Wait()
}

func (h *HealthChecker) CheckOnce(ctx context.Context) {
	h.checkAll(ctx)
}

func (h *HealthChecker) CheckHealth(ctx context.Context) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.checkAll(ctx)
		}
	}
}
