package app

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestCustomHandler_Handle(t *testing.T) {
	writer := &bytes.Buffer{}
	handler := &customHandler{
		writer: writer,
		level:  slog.LevelInfo,
	}

	tests := []struct {
		name        string
		recordSetup func() slog.Record
		expectedLog string
		shouldLog   bool
	}{
		{
			name: "INFO log level",
			recordSetup: func() slog.Record {
				record := slog.NewRecord(time.Now(), slog.LevelInfo, "Test INFO log", 0)
				return record
			},
			expectedLog: "INFO[",
			shouldLog:   true,
		},
		{
			name: "DEBUG log level (below threshold)",
			recordSetup: func() slog.Record {
				record := slog.NewRecord(time.Now(), slog.LevelDebug, "Test DEBUG log", 0)
				return record
			},
			expectedLog: "",
			shouldLog:   false,
		},
		{
			name: "Log with attributes",
			recordSetup: func() slog.Record {
				record := slog.NewRecord(time.Now(), slog.LevelInfo, "Log with attributes", 0)
				record.AddAttrs(slog.String("key1", "value1"), slog.String("key2", "value2"))
				return record
			},
			expectedLog: "INFO[",
			shouldLog:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer.Reset()
			record := test.recordSetup()
			err := handler.Handle(context.Background(), record)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			logOutput := writer.String()
			if test.shouldLog && !contains(logOutput, test.expectedLog) {
				t.Errorf("expected log output to contain %q, got %q", test.expectedLog, logOutput)
			}
			if !test.shouldLog && logOutput != "" {
				t.Errorf("expected no log output, got %q", logOutput)
			}
		})
	}
}

func TestCustomHandler_Enabled(t *testing.T) {
	handler := &customHandler{
		level: slog.LevelWarn,
	}

	tests := []struct {
		name      string
		level     slog.Level
		isEnabled bool
	}{
		{"Below threshold", slog.LevelInfo, false},
		{"At threshold", slog.LevelWarn, true},
		{"Above threshold", slog.LevelError, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if handler.Enabled(context.Background(), test.level) != test.isEnabled {
				t.Errorf("expected Enabled(%v) = %v", test.level, test.isEnabled)
			}
		})
	}
}

func TestCustomHandler_WithAttrs(t *testing.T) {
	handler := &customHandler{}
	newHandler := handler.WithAttrs(nil)

	if newHandler == nil {
		t.Error("expected WithAttrs to return a non-nil handler")
	}
}

func TestCustomHandler_WithGroup(t *testing.T) {
	handler := &customHandler{}
	newHandler := handler.WithGroup("testGroup")

	if newHandler == nil {
		t.Error("expected WithGroup to return a non-nil handler")
	}
}

func TestNewCustomLogger(t *testing.T) {
	os.Setenv("TB_LOGLEVEL", "debug")
	defer os.Unsetenv("TB_LOGLEVEL")

	logger, level := NewCustomLogger()
	if logger == nil {
		t.Error("expected NewCustomLogger to return a valid logger")
	}
	if level != slog.LevelDebug {
		t.Errorf("expected log level to be DEBUG, got %v", level)
	}

	os.Unsetenv("TB_LOGLEVEL")
	_, level = NewCustomLogger()
	if level != slog.LevelInfo {
		t.Errorf("expected log level to default to INFO, got %v", level)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 && bytes.Contains([]byte(haystack), []byte(needle))
}
