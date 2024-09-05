package main

import (
	"fmt"
	"log/slog"
	"os"
)

type App struct {
	Logger *slog.Logger
	Config *Config
}

func NewApp() (*App, error) {
	a := App{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}

	if config, err := establishConfig(); err != nil {
		return nil, err
	} else {
		a.Config = config
		return &a, nil
	}
}

func (a *App) ParseArgs() bool {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("v3.0.3-2024-09-01-23-18-10")
		return false
	}
	return true
}
