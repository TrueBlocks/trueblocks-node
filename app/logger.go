package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type CustomHandler struct {
	writer io.Writer
	level  slog.Level
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < h.level {
		return nil
	}
	levels := map[string]string{
		"DEBU": "DEBG",
		"INFO": "INFO",
		"WARN": "WARN",
		"ERRO": "EROR",
	}
	lev := r.Level.String()[:4]
	timeFormat := r.Time.Format("02-01|15:04:05.000")

	// Format the message to be exactly n characters wide
	formattedMessage := fmt.Sprintf("%-18.18s", r.Message)

	logMsg := fmt.Sprintf("%4.4s[%s] %s ", levels[lev], timeFormat, formattedMessage)
	r.Attrs(func(attr slog.Attr) bool {
		logMsg += fmt.Sprintf(" %s=%v", colors.Green+attr.Key+colors.Off, attr.Value)
		return true
	})
	fmt.Fprintln(h.writer, logMsg)
	return nil
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}

func NewCustomLogger() (*slog.Logger, slog.Level) {
	logger.SetLoggerWriter(io.Discard)
	logLevel := slog.LevelInfo
	if ll, ok := os.LookupEnv("TB_LOGLEVEL"); ok {
		switch strings.ToLower(ll) {
		case "debug":
			logLevel = slog.LevelDebug
		case "info":
			logLevel = slog.LevelInfo
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		}
	}
	customHandler := &CustomHandler{
		writer: os.Stderr,
		level:  logLevel,
	}
	return slog.New(customHandler), logLevel
}
