package main

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type Feature string

const (
	Scrape  Feature = "scrape"
	Monitor Feature = "monitor"
	Api     Feature = "api"
)

type InitMode string

const (
	All    InitMode = "all"
	Blooms InitMode = "blooms"
	None   InitMode = "none"
)

type OnOff string

const (
	On  OnOff = "on"
	Off OnOff = "off"
)

type App struct {
	Logger   *slog.Logger
	Config   Config
	InitMode InitMode
	Monitor  OnOff
	Api      OnOff
	Sleep    int
	Busy     bool
}

func NewApp() *App {
	logger.SetLoggerWriter(io.Discard) // we never want core to log anything
	logLevel := slog.LevelInfo
	if ll, ok := os.LookupEnv("TB_NODE_LOGLEVEL"); ok {
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
	opts := slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("15:04:05"))
			}
			return a
		},
	}

	app := &App{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, &opts)),
		Sleep:  4,
		Config: Config{
			ProviderMap: make(map[string]string),
		},
		Api:      On,
		Monitor:  Off,
		InitMode: Blooms,
	}
	// app.Logger.Info("Starting", "logLevel", logLevel)

	return app
}

func (a *App) IsOn(feature Feature) bool {
	switch feature {
	case Scrape:
		return a.InitMode == All || a.InitMode == Blooms
	case Monitor:
		return a.Monitor == On
	case Api:
		return a.Api == On
	}
	return false
}
