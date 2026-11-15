package loadbalancer

import "testing"

func TestLoadbalancing(t *testing.T) {
	tests := []struct {
		key     string
		servers []string
		want    string
	}{
		{key: "/api/test", servers: []string{"A", "B", "C"}, want: "A"},
		{key: "/api/test", servers: []string{"A", "B", "C"}, want: "B"},
		{key: "/api/test", servers: []string{"A", "B", "C"}, want: "C"},
		{key: "/api/test", servers: []string{"A", "B", "C"}, want: "A"},
	}

	lb := NewLoadBalancer()

	for _, test := range tests {
		got := lb.NextServer(test.key, test.servers)

		if got != test.want {
			t.Errorf("expected %v, got %v", test.want, got)
		}
	}
}

func TestLoadBalancerWithEmptyServers(t *testing.T) {
	lb := NewLoadBalancer()

	got := lb.NextServer("/api/test", []string{})

	if got != "" {
		t.Errorf("expected empty string, got %v", got)
	}
}
