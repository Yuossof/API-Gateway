package gatewayconfig

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Yuossof/go-api-gateway/internal/cluster"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Route struct {
	Path    string `yaml:"path" validate:"required"`
	Cluster string `yaml:"cluster" validate:"required"`
	Target  string `yaml:"target" validate:"required"`
}

type RateLimit struct {
	MaxTimeForToken int `yaml:"max_time_for_token" validate:"required,gt=0"`
	NumberOfTokens  int `yaml:"number_of_tokens" validate:"required,gt=0"`
}

type Config struct {
	Clusters            []cluster.Cluster           `yaml:"clusters" validate:"required,min=1,dive"`
	Routes              []Route                     `yaml:"routes" validate:"required,min=1,dive"`
	RateLimit           RateLimit                   `yaml:"ratelimit" validate:"required"`
	AllServersByCluster map[string][]cluster.Server `yaml:"-"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	validate := validator.New()

	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.AllServersByCluster = make(map[string][]cluster.Server)

	for _, currentCluster := range config.Clusters {
		config.AllServersByCluster[currentCluster.Name] = currentCluster.Servers
	}

	slog.Info("config loaded successfully", "path", filePath)
	return &config, nil
}
