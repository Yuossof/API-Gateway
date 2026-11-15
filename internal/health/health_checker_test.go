package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Yuossof/go-api-gateway/internal/cluster"
)

func TestGetHealthyServers(t *testing.T) {
	checker := NewHealthChecker(nil)
	checker.status["http://healthy"] = true
	checker.status["http://unhealthy"] = false

	servers := checker.GetHealthyServers([]cluster.Server{
		{URL: "http://healthy"},
		{URL: "http://unhealthy"},
	})

	want := []cluster.Server{{URL: "http://healthy"}}
	if !reflect.DeepEqual(want, servers) {
		t.Errorf("expected %v, got %v", want, servers)
	}
}

func TestCheckOnce(t *testing.T) {
	healthyServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/health" {
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	))
	defer healthyServer.Close()

	unhealthyServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer unhealthyServer.Close()

	servers := []cluster.Server{
		{URL: healthyServer.URL},
		{URL: unhealthyServer.URL},
	}
	checker := NewHealthChecker(map[string][]cluster.Server{
		"test": servers,
	})

	checker.CheckOnce(context.Background())

	got := checker.GetHealthyServers(servers)
	want := []cluster.Server{{URL: healthyServer.URL}}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestSendReq(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{name: "OK", status: http.StatusOK, expected: true},
		{name: "Created", status: http.StatusCreated, expected: true},
		{name: "Internal Server Error", status: http.StatusInternalServerError, expected: false},
		{name: "Not Found", status: http.StatusNotFound, expected: false},
	}

	checker := NewHealthChecker(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.status)
				},
			))
			defer server.Close()

			if got := checker.SendReq(context.Background(), server.URL); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
