package configmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func (m *Manager) resolveConfigPath(input string) (string, error) {
	cleaned := filepath.Clean(input)

	// لو المستخدم بعت اسم ملف بس بدون مسار (مثلاً v1.yaml) وكان فيه versionsDir محدد
	if !strings.Contains(cleaned, string(filepath.Separator)) && m.versionsDir != "" {
		cleaned = filepath.Join(m.versionsDir, cleaned)
	}

	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(absPath))
	if ext != ".yaml" && ext != ".yml" {
		return "", fmt.Errorf("unsupported config extension: %s", ext)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("config path is a directory: %s", absPath)
	}

	return absPath, nil
}

func (m *Manager) Activate(configPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	resolvedPath, err := m.resolveConfigPath(configPath)
	if err != nil {
		return err
	}

	gatewayConfig, err := gatewayconfig.LoadConfig(resolvedPath)
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

	m.activeVersion = resolvedPath
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
			Path string `json:"path"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil || strings.TrimSpace(request.Path) == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := m.Activate(request.Path); err != nil {
			http.Error(w, "activation failed: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success":        true,
			"active_version": m.activeVersion,
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
