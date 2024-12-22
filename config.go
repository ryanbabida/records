package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	defaultFilePath = "somedefault.json"
	defaultPort     = "8080"
)

type Config struct {
	DataFilePath    *string `json:"dataFilePath"`
	Port            *string `json:"port"`
	ArtworkFilePath *string `json:"artworkFilePath"`
	AudioFilePath   *string `json:"audioFilePath"`
}

func NewConfig(filepath string, opts ...func(c *Config)) (*Config, error) {
	configFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file at '%s': %w", filepath, err)
	}

	b, err := io.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	filePath := defaultFilePath
	port := defaultPort
	cfg := &Config{DataFilePath: &filePath, Port: &port}

	err = json.Unmarshal(b, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	for _, o := range opts {
		o(cfg)
	}

	return cfg, nil
}
