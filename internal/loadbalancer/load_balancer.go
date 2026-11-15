package loadbalancer

import "sync"

type LoadBalancer struct {
	mu      sync.Mutex
	indexes map[string]int
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		indexes: make(map[string]int),
	}
}

func (lb *LoadBalancer) NextServer(key string, servers []string) string {
	if len(servers) == 0 {
		return ""
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	currentIndex := lb.indexes[key]
	server := servers[currentIndex]

	lb.indexes[key] = (currentIndex + 1) % len(servers)

	return server
	// server := servers[]
}
