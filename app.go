package main

import (
	"log/slog"
)

type App struct {
	Logger *slog.Logger
	Config Config
	Busy   bool
}
