package main

import (
	"os"
	"sync"
)

func main() {
	a := NewApp()
	if ok := a.parseArgs(); !ok {
		return
	}

	// Establish the configuration file
	if err := a.establishConfig(); err != nil {
		a.Logger.Error(err.Error())

	} else {
		os.Exit(0)

		// Start the API server. Don't wait for it to return, it never does.
		if a.Api.IsOn() {
			a.serve()
		}

		// Start forever loop to scrape and (optionally) monitor the chain
		wg := sync.WaitGroup{}
		wg.Add(1)
		go a.scrape(&wg)
		if a.Monitor.IsOn() {
			wg.Add(1)
			go a.monitor(&wg)
		}
		wg.Wait()
	}
}
