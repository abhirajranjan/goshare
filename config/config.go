package config

import (
	"flag"
	"goshare/resources"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./config.yaml", "set config")
	flag.Parse()
}

func LoadConfig() (resources.Config, error) {
	var cfg resources.Config
	cfg.LoadDefault()

	fd, err := os.Open(configPath)
	if err != nil {
		slog.Warn("loading default config")
		return cfg, nil
	}

	if err := yaml.NewDecoder(fd).Decode(&cfg); err != nil {
		return cfg, errors.Wrap(err, "yaml.Decode")
	}

	return cfg, nil
}
