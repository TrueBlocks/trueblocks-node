package main

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func main() {
	if ok := parseArgs(); !ok {
		return
	}

	// Establish the app with a logger (and turn off the core's logging)
	a := App{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	logger.SetLoggerWriter(io.Discard)

	// Establish the configuration file
	if err := a.establishConfig(); err != nil {
		a.Logger.Error(err.Error())

	} else {
		// Start the API server. Don't wait for it to finish. It never does.
		a.serve()

		// Start two forever loops to scrape and monitor the chain
		wg := sync.WaitGroup{}
		wg.Add(2)
		go a.scrape(&wg)
		go a.monitor(&wg)
		wg.Wait()
	}
}
