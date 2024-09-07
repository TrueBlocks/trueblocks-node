package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	a := NewApp()
	if cont, err := a.parseArgs(); !cont {
		return // don't continue
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// Establish the configuration file
	if err := a.establishConfig(); err != nil {
		a.Logger.Error(err.Error())

	} else {
		a.Logger.Info("Starting trueBlocks-node with...", "api", a.IsOn(Api), "scrape", a.IsOn(Scrape), "monitor", a.IsOn(Monitor))

		// Start the API server. It runs in its own goroutine.
		if a.IsOn(Api) {
			a.serve()
		}

		// Start forever loop to scrape and (optionally) monitor the chain
		wg := sync.WaitGroup{}

		if a.IsOn(Scrape) {
			wg.Add(1)
			go a.scrape(&wg)
		}

		if a.IsOn(Monitor) {
			wg.Add(1)
			go a.monitor(&wg)
		}

		wg.Wait()
	}
}
