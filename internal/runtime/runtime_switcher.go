package runtime

import (
	"net/http"
	"sync"
)

type RuntimeSwitcher struct {
	runtime *Runtime
	mu      sync.RWMutex
}

func NewRuntimeSwitcher() *RuntimeSwitcher {
	return &RuntimeSwitcher{}
}

func (rs *RuntimeSwitcher) SetRuntime(rt *Runtime) {
	rs.mu.Lock()
	rs.runtime = rt
	rs.mu.Unlock()
}

func (rs *RuntimeSwitcher) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	rs.mu.RLock()
	current := rs.runtime
	rs.mu.RUnlock()

	current.handler.ServeHTTP(w, r)
}
