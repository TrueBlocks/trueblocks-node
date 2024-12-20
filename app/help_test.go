package app

import (
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
		api        OnOff
		monitor    OnOff
		sleep      int
	}{
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
			expectErr:  "", // helpText,
		},
		{
			name:       "Unknown Flag",
			args:       []string{"--junk"},
			expectBool: true,
			expectErr:  "unknown option:--junk\n" + helpText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			os.Args = append([]string{"test"}, tt.args...)
			resultBool, err := app.ParseArgs()

			if resultBool != tt.expectBool {
				t.Errorf("expected %v, got %v", tt.expectBool, resultBool)
			}

			if (err == nil && tt.expectErr != "") || (err != nil && err.Error() != tt.expectErr) {
				t.Errorf("expected error %q, got %q", tt.expectErr, err)
			}

			if tt.initMode != "" && app.InitMode != tt.initMode {
				t.Errorf("expected InitMode %v, got %v", tt.initMode, app.InitMode)
			}

			if tt.api != "" && app.Api != tt.api {
				t.Errorf("expected Api %v, got %v", tt.api, app.Api)
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
