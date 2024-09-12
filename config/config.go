package config

import (
	"encoding/json"
	"path/filepath"
)

// Config carries various configuration related settings
type Config struct {
	ConfigPath  string            `json:"configPath"`
	ProviderMap map[string]string `json:"providers"` // chain to provider
	Targets     []string          `json:"targets"`
}

// String returns the string representation of the Config type
func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// CachePath returns the path to the TrueBlocks cache directory
func (c *Config) CachePath() string {
	return filepath.Join(c.ConfigPath, "cache")
}

// Index returns the path to the TrueBlocks Unchained Index directory
func (c *Config) IndexPath() string {
	return filepath.Join(c.ConfigPath, "unchained")
}
