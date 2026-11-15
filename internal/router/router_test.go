package router

import (
	"reflect"
	"testing"

	"github.com/Yuossof/go-api-gateway/internal/gatewayconfig"
)

func TestMatchRoute(t *testing.T) {
	route := []gatewayconfig.Route{
		{
			Path:    "/api/test",
			Cluster: "",
			Target:  "/api/test",
		},
		{
			Path:    "/api/test2",
			Cluster: "",
			Target:  "/api/test2",
		},
	}

	r, _ := MatchRoute("/api/test", route)

	expected := gatewayconfig.Route{
		Path:    "/api/test",
		Cluster: "",
		Target:  "/api/test",
	}

	if !reflect.DeepEqual(expected, r) {
		t.Errorf("expected %v, got %v", expected, r)
	}
}
