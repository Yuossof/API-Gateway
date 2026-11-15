package proxy

import (
	"testing"

	"github.com/Yuossof/go-api-gateway/internal/cluster"
)

func TestLoadProxies(t *testing.T) {
	clusters := []cluster.Cluster{
		{
			Name: "c1",
			Servers: []cluster.Server{
				{URL: "http://localhost:8001"},
				{URL: "http://localhost:8002"},
			},
		},

		{
			Name: "c2",
			Servers: []cluster.Server{
				{URL: "http://localhost:9001"},
				{URL: "http://localhost:9002"},
			},
		},
	}

	proxies, _ := LoadProxies(clusters)

	expectedURLs := []string{
		"http://localhost:8001",
		"http://localhost:8002",
		"http://localhost:9001",
		"http://localhost:9002",
	}

	if len(proxies) != len(expectedURLs) {
		t.Fatalf("expected %d proxies, got %d", len(expectedURLs), len(proxies))
	}

	for _, url := range expectedURLs {
		proxy, exists := proxies[url]

		if !exists {
			t.Errorf("expected proxy for URL %s", url)
			continue
		}

		if proxy == nil {
			t.Errorf("expected non-nil proxy for URL %s", url)
		}
	}
}
