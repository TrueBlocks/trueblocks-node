package config

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestConfigMethods(t *testing.T) {
	config := Config{
		ConfigPath:  "/path/to/config",
		ProviderMap: map[string]string{"chain1": "provider1", "chain2": "provider2"},
		Targets:     []string{"target1", "target2"},
	}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name: "String",
			method: func() string {
				bytes, _ := json.Marshal(config)
				return string(bytes)
			},
			expected: config.String(),
		},
		{
			name: "CachePath",
			method: func() string {
				return filepath.Join(config.ConfigPath, "cache")
			},
			expected: config.CachePath(),
		},
		{
			name: "IndexPath",
			method: func() string {
				return filepath.Join(config.ConfigPath, "unchained")
			},
			expected: config.IndexPath(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.method()
			if actual != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, actual)
			}
		})
	}
}
