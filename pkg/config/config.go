package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Provider string       `json:"provider"`
	Gemini   GeminiConfig `json:"gemini"`
}

type GeminiConfig struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Provider: "gemini",
		Gemini: GeminiConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
			Model:  "gemini-2.5-flash",
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}

	configPath := filepath.Join(home, ".commitgen.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}

	if envKey := os.Getenv("GEMINI_API_KEY"); envKey != "" {
		cfg.Gemini.APIKey = envKey
	}

	return cfg, nil
}
