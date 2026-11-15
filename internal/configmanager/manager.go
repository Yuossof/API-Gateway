package configmanager

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Yuossof/go-api-gateway/internal/gatewayconfig"
	"github.com/Yuossof/go-api-gateway/internal/runtime"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	mu            sync.Mutex
	versionsDir   string
	activeVersion string
	redis         *redis.Client
	runtime       *runtime.Runtime
	switcher      *runtime.RuntimeSwitcher
}

func NewManager(versionsDir string, redis *redis.Client, switcher *runtime.RuntimeSwitcher) *Manager {
	return &Manager{
		versionsDir: versionsDir,
		redis:       redis,
		switcher:    switcher,
	}
}

func (m *Manager) Activate(versionName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	gatewayConfig, err := gatewayconfig.LoadConfig(versionName, m.versionsDir)

	if err != nil {
		return err
	}

	newRuntime, err := runtime.BuildRuntime(gatewayConfig, m.redis)

	if err != nil {
		return err
	}

	previousRuntime := m.runtime
	m.runtime = newRuntime

	m.switcher.SetRuntime(newRuntime)

	m.activeVersion = versionName
	if previousRuntime != nil {
		previousRuntime.Stop()
	}

	return nil
}

func (m *Manager) ActivationListener() error {
	const socketPath = "/tmp/gateway.sock"

	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	defer os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer listener.Close()

	if err := os.Chmod(socketPath, 0600); err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/activate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)

		var request struct {
			Version string `json:"version"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if request.Version == "" ||
			filepath.Base(request.Version) != request.Version ||
			filepath.Ext(request.Version) != ".yaml" {
			http.Error(w, "invalid version", http.StatusBadRequest)
			return
		}

		if err := m.Activate(request.Version); err != nil {
			http.Error(w, "activation failed", http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success":        true,
			"active_version": request.Version,
		})
	})

	server := &http.Server{
		Handler: mux,
	}

	err = server.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
