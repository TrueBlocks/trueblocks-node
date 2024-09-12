package app

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-node/v3/config"
)

// Feature is a type that represents the features of the app
type Feature string

const (
	// Scrape represents the scraper feature. The scraper may not be disabled.
	Scrape Feature = "scrape"
	// Monitor represents the monitor feature. The monitor is Off by default. Enable
	// it with the `--monitor on` option.
	Monitor Feature = "monitor"
	// Api represents the API feature. The api is On by default. Disable it
	// with the `--api off` option.
	Api Feature = "api"
)

// InitMode is a type that represents the initialization for the Unchained Index. It
// applies to the `--init` option.
type InitMode string

const (
	// All cause the initialization to download both the bloom filters and the index
	// portions of the Unchained Index.
	All InitMode = "all"
	// Blooms cause the initialization to download only the bloom filters portion of
	// the Unchained Index.
	Blooms InitMode = "blooms"
	// None cause the app to not download any part of the Unchained Index. It will be
	// built from scratch with the scraper.
	None InitMode = "none"
)

// OnOff is a type that represents a boolean value that can be either "on" or "off".
type OnOff string

const (
	// On is the "on" value for a feature. It applies to the `--monitor` and `--api` options.
	On OnOff = "on"
	// Off is the "off" value for a feature. It applies to the `--monitor` and `--api` options.
	Off OnOff = "off"
)

// App is the main struct for the app. It contains the logger, the configuration, and the
// state of the app.
type App struct {
	Logger   *slog.Logger
	Config   config.Config
	InitMode InitMode
	Monitor  OnOff
	Api      OnOff
	Sleep    int
}

// NewApp creates a new App instance with the default values.
func NewApp() *App {
	logger.SetLoggerWriter(io.Discard) // we never want core to log anything
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
		Config: config.Config{
			ProviderMap: make(map[string]string),
		},
		Api:      On,
		Monitor:  Off,
		InitMode: Blooms,
	}
	// app.Logger.Info("Starting", "logLevel", logLevel)

	return app
}

// IsOn returns true if the given feauture is enabled. It returns false otherwise.
func (a *App) IsOn(feature Feature) bool {
	switch feature {
	case Scrape:
		return a.InitMode != None
	case Monitor:
		return a.Monitor == On
	case Api:
		return a.Api == On
	}
	return false
}
