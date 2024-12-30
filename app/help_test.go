package app

import (
	"fmt"
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectBool bool
		expectErr  string
		initMode   InitMode
		scrape     OnOff
		api        OnOff
		ipfs       OnOff
		monitor    OnOff
		sleep      int
	}{
		// ------------------------------
		{
			name:       "No Arguments",
			args:       []string{},
			expectBool: true,
			expectErr:  "",
		},
		{
			name:       "Valid Init",
			args:       []string{"--init", "all"},
			expectBool: true,
			expectErr:  "",
			initMode:   All,
		},
		{
			name:       "Missing Init Argument",
			args:       []string{"--init"},
			expectBool: true,
			expectErr:  "missing argument for --init",
		},
		{
			name:       "Invalid Init Argument",
			args:       []string{"--init", "invalid"},
			expectBool: true,
			expectErr:  "parsing --init: invalid value for mode: invalid",
		},

		// ------------------------------
		{
			name:       "Valid Scrape",
			args:       []string{"--scrape", "on"},
			expectBool: true,
			expectErr:  "",
			scrape:     On,
		},
		{
			name:       "Missing Scrape Argument",
			args:       []string{"--scrape"},
			expectBool: true,
			expectErr:  "missing argument for --scrape",
		},
		{
			name:       "Invalid Scrape Argument",
			args:       []string{"--scrape", "invalid"},
			expectBool: true,
			expectErr:  "parsing --scrape: invalid value for onOff: invalid",
		},

		// ------------------------------
		{
			name:       "Valid API",
			args:       []string{"--api", "on"},
			expectBool: true,
			expectErr:  "",
			api:        On,
		},
		{
			name:       "Missing API Argument",
			args:       []string{"--api"},
			expectBool: true,
			expectErr:  "missing argument for --api",
		},
		{
			name:       "Invalid API Argument",
			args:       []string{"--api", "invalid"},
			expectBool: true,
			expectErr:  "parsing --api: invalid value for onOff: invalid",
		},

		// ------------------------------
		{
			name:       "Valid IPFS",
			args:       []string{"--ipfs", "on"},
			expectBool: true,
			expectErr:  "",
			ipfs:       On,
		},
		{
			name:       "Missing IPFS Argument",
			args:       []string{"--ipfs"},
			expectBool: true,
			expectErr:  "missing argument for --ipfs",
		},
		{
			name:       "Invalid IPFS Argument",
			args:       []string{"--ipfs", "invalid"},
			expectBool: true,
			expectErr:  "parsing --ipfs: invalid value for onOff: invalid",
		},

		// ------------------------------
		{
			name:       "Valid Monitor",
			args:       []string{"--monitor", "off"},
			expectBool: true,
			expectErr:  "",
			monitor:    Off,
		},
		{
			name:       "Missing Monitor Argument",
			args:       []string{"--monitor"},
			expectBool: true,
			expectErr:  "missing argument for --monitor",
		},
		{
			name:       "Invalid Monitor Argument",
			args:       []string{"--monitor", "invalid"},
			expectBool: true,
			expectErr:  "parsing --monitor: invalid value for onOff: invalid",
		},

		// ------------------------------
		{
			name:       "Valid Sleep",
			args:       []string{"--sleep", "60"},
			expectBool: true,
			expectErr:  "",
			sleep:      60,
		},
		{
			name:       "Missing Sleep Argument",
			args:       []string{"--sleep"},
			expectBool: true,
			expectErr:  "missing argument for --sleep",
		},
		{
			name:       "Invalid Sleep Argument",
			args:       []string{"--sleep", "invalid"},
			expectBool: true,
			expectErr:  "parsing --sleep: invalid value for sleep: invalid",
		},

		// ------------------------------
		{
			name:       "Version Flag",
			args:       []string{"--version"},
			expectBool: false,
			expectErr:  "",
		},
		{
			name:       "Help Flag",
			args:       []string{"--help"},
			expectBool: false,
			expectErr:  "",
		},
		{
			name:       "Unknown Flag",
			args:       []string{"--junk"},
			expectBool: true,
			expectErr:  "unknown option:--junk\n" + helpText,
		},
	}

	os.Setenv("TEST_MODE", "true")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			os.Args = append([]string{"test"}, tt.args...)
			resultBool, _, err := app.ParseArgs()
			fmt.Println(tt.name, tt.args, resultBool, err)

			if resultBool != tt.expectBool {
				t.Errorf("expected %v, got %v", tt.expectBool, resultBool)
			}

			if (err == nil && tt.expectErr != "") || (err != nil && err.Error() != tt.expectErr) {
				t.Errorf("expected error %q, got %q", tt.expectErr, err)
			}

			if tt.initMode != "" && app.InitMode != tt.initMode {
				t.Errorf("expected InitMode %v, got %v", tt.initMode, app.InitMode)
			}

			if tt.scrape != "" && app.Scrape != tt.scrape {
				t.Errorf("expected Scraper %v, got %v", tt.scrape, app.Scrape)
			}

			if tt.api != "" && app.Api != tt.api {
				t.Errorf("expected Api %v, got %v", tt.api, app.Api)
			}

			if tt.ipfs != "" && app.Ipfs != tt.ipfs {
				t.Errorf("expected Ipfs %v, got %v", tt.ipfs, app.Ipfs)
			}

			if tt.monitor != "" && app.Monitor != tt.monitor {
				t.Errorf("expected Monitor %v, got %v", tt.monitor, app.Monitor)
			}

			if tt.sleep != 0 && app.Sleep != tt.sleep {
				t.Errorf("expected Sleep %v, got %v", tt.sleep, app.Sleep)
			}
		})
	}
}
