package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/TrueBlocks/trueblocks-node/v4/app"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
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
		a.Fatal(err)

	} else {
		a.Logger.Info("trueblocks-node", "scrape", a.State(app.Scrape), "api", a.State(app.Api), "monitor", a.State(app.Monitor), "init-mode", a.InitMode)

		// Start the API server. It runs in its own goroutine.
		var apiUrl string
		if a.IsOn(app.Api) {
			go sdk.StartApiServer(&apiUrl)
		}

		// Start forever loop to scrape and (optionally) monitor the chain
		wg := sync.WaitGroup{}

		if a.IsOn(app.Scrape) {
			wg.Add(1)
			a.Logger.Info("start scraper...")
			go a.RunScraper(&wg)
			a.Logger.Info("scraper started...")
		}

		if a.IsOn(app.Monitor) {
			wg.Add(1)
			a.Logger.Info("start monitors...")
			go a.RunMonitor(&wg)
			a.Logger.Info("monitors started...")
		}

		if len(apiUrl) > 0 {
			a.Logger.Info("api is runing", "apiUrl", apiUrl)
		}

		wg.Wait()
	}
}
