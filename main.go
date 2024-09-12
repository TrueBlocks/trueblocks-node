package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/TrueBlocks/trueblocks-node/v3/app"
)

func main() {
	a := app.NewApp()
	if cont, err := a.ParseArgs(); !cont {
		return // don't continue
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// Establish the configuration file
	if err := a.EstablishConfig(); err != nil {
		a.Logger.Error(err.Error())

	} else {
		a.Logger.Info("Starting trueBlocks-node with...", "api", a.IsOn(app.Api), "init", a.InitMode, "monitor", a.IsOn(app.Monitor))

		// Start the API server. It runs in its own goroutine.
		if a.IsOn(app.Api) {
			a.Logger.Info("Starting Api server...")
			a.RunServer()
			a.Logger.Info("Api server started...")
		}

		// Start forever loop to scrape and (optionally) monitor the chain
		wg := sync.WaitGroup{}

		if a.IsOn(app.Scrape) {
			wg.Add(1)
			a.Logger.Info("Starting scraper...")
			go a.RunScraper(&wg)
			a.Logger.Info("Scraper started...")
		}

		if a.IsOn(app.Monitor) {
			wg.Add(1)
			a.Logger.Info("Starting monitors...")
			go a.RunMonitor(&wg)
			a.Logger.Info("Monitors started...")
		}

		wg.Wait()
	}
}
