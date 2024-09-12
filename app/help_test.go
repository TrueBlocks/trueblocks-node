package app

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectBool bool
		expectErr  error
		initMode   InitMode
		api        OnOff
		monitor    OnOff
		sleep      int
	}{
		{
			name:       "No Arguments",
			args:       []string{},
			expectBool: true,
			expectErr:  nil,
		},
		{
			name:       "Valid Init",
			args:       []string{"--init", "all"},
			expectBool: true,
			expectErr:  nil,
			initMode:   All,
		},
		{
			name:       "Missing Init Argument",
			args:       []string{"--init"},
			expectBool: true,
			expectErr:  fmt.Errorf("missing argument for --init"),
		},
		{
			name:       "Invalid Init Argument",
			args:       []string{"--init", "invalid"},
			expectBool: true,
			expectErr:  fmt.Errorf("invalid value: invalid"),
		},
		{
			name:       "Valid API",
			args:       []string{"--api", "on"},
			expectBool: true,
			expectErr:  nil,
			api:        On,
		},
		{
			name:       "Missing API Argument",
			args:       []string{"--api"},
			expectBool: true,
			expectErr:  fmt.Errorf("missing argument for --init"), // Correct this to match the actual error message
		},
		{
			name:       "Invalid API Argument",
			args:       []string{"--api", "invalid"},
			expectBool: true,
			expectErr:  fmt.Errorf("invalid value: invalid"),
		},
		{
			name:       "Valid Monitor",
			args:       []string{"--monitor", "off"},
			expectBool: true,
			expectErr:  nil,
			monitor:    Off,
		},
		{
			name:       "Missing Monitor Argument",
			args:       []string{"--monitor"},
			expectBool: true,
			expectErr:  fmt.Errorf("missing argument for --monitor"),
		},
		{
			name:       "Invalid Monitor Argument",
			args:       []string{"--monitor", "invalid"},
			expectBool: true,
			expectErr:  fmt.Errorf("invalid value: invalid"),
		},
		{
			name:       "Valid Sleep",
			args:       []string{"--sleep", "60"},
			expectBool: true,
			expectErr:  nil,
			sleep:      int(60 * time.Second),
		},
		{
			name:       "Missing Sleep Argument",
			args:       []string{"--sleep"},
			expectBool: true,
			expectErr:  fmt.Errorf("missing argument for --sleep"),
		},
		{
			name:       "Invalid Sleep Argument",
			args:       []string{"--sleep", "invalid"},
			expectBool: true,
			expectErr:  fmt.Errorf("invalid value for --sleep: invalid"),
		},
		{
			name:       "Version Flag",
			args:       []string{"--version"},
			expectBool: false,
			expectErr:  nil,
		},
		{
			name:       "Help Flag",
			args:       []string{"--help"},
			expectBool: false,
			expectErr:  fmt.Errorf(helpText),
		},
		// {
		// 	name:       "Unknown Flag",
		// 	args:       []string{"--junk"},
		// 	expectBool: false,
		// 	expectErr:  fmt.Errorf("%s%s", "Unknown option:--junk\n", helpText),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{}
			os.Args = append([]string{"test"}, tt.args...)
			resultBool, err := app.ParseArgs()
			if resultBool != tt.expectBool {
				t.Errorf("expected %v, got %v", tt.expectBool, resultBool)
			}
			if (err == nil) != (tt.expectErr == nil) || (err != nil && err.Error() != tt.expectErr.Error()) {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
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
