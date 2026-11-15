package proxy

import (
	"net/http"

	"github.com/Yuossof/go-api-gateway/internal/cluster"
)

type Proxies map[string]http.Handler

func LoadProxies(clusters []cluster.Cluster) (Proxies, error) {
	proxies := make(Proxies)

	for _, u := range clusters {

		for _, s := range u.Servers {
			proxy, err := NewProxy(s.URL)
			if err != nil {
				return nil, err
			}

			proxies[s.URL] = proxy
		}

	}

	return proxies, nil
}
