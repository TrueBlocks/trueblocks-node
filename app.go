package main

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type SyncType string

const (
	All    SyncType = "all"
	Blooms SyncType = "blooms"
)

type OnOff string

const (
	On  OnOff = "on"
	Off OnOff = "off"
)

func (o OnOff) IsOn() bool {
	return o == On
}

type App struct {
	Logger  *slog.Logger
	Config  Config
	Mode    SyncType
	Monitor OnOff
	Api     OnOff
	Sleep   time.Duration
	Busy    bool
}

func NewApp() *App {
	logger.SetLoggerWriter(io.Discard) // we never want core to log anything
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug, // Set the minimum log level
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("15:04:05"))
			}
			return a
		},
	}

	return &App{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, &opts)),
		Sleep:  4,
		Config: Config{
			ProviderMap: make(map[string]string),
		},
	}
}
