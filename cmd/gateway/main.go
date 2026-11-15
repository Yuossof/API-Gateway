package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yuossof/go-api-gateway/internal/config"
	"github.com/Yuossof/go-api-gateway/internal/configmanager"
	"github.com/Yuossof/go-api-gateway/internal/matrics"
	"github.com/Yuossof/go-api-gateway/internal/redisclient"
	"github.com/Yuossof/go-api-gateway/internal/runtime"
)

func main() {

	ctx := context.Background()

	appConfig, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create redis client
	redis, err := redisclient.New(ctx, appConfig.RedisURL)

	if err != nil {
		slog.Error("failed to create redis client", "error", err)
		os.Exit(1)
	}

	switcher := runtime.NewRuntimeSwitcher()

	manager := configmanager.NewManager("configs/gateway/versions/", redis, switcher)

	err = manager.Activate("v1.yaml") // TODO: make the activation happen through CLI

	if err != nil {
		slog.Error("activation error", "error", err)
		os.Exit(1)
	}

	activationError := make(chan error, 1)
	go func() {
		if err := manager.ActivationListener(); err != nil {
			activationError <- err
		}
	}()

	matrics.StartMatricsServer()

	mux := http.NewServeMux()

	mux.Handle("/", switcher)

	server := &http.Server{
		Addr:         ":" + appConfig.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("server is running on :%s", appConfig.Port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case err := <-activationError:
		log.Printf("Activation listener error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		return
	}

	slog.Info("Server exited properly")
}
