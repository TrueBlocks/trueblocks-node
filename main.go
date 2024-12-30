package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/TrueBlocks/trueblocks-node/v4/app"
	"github.com/TrueBlocks/trueblocks-sdk/v4/services"
)

func main() {
	a := app.NewApp()
	if cont, err := a.ParseArgs(); !cont {
		return // don't continue
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if err := a.EstablishConfig(); err != nil {
		a.Fatal(err)
	} else {
		a.Logger.Info("trueblocks-node", "scrape", a.State(app.Scrape), "api", a.State(app.Api), "ipfs", a.State(app.Ipfs), "monitor", a.State(app.Monitor), "init-mode", a.InitMode)

		cleanupChan := make(chan string, 4)
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt)

		runningServices := 0
		if a.IsOn(app.Scrape) {
			runningServices++
			go func() {
				a.Logger.Info("scraper started...")
				a.RunScraper(nil)
				a.Logger.Info("scraper cleanup completed")
				cleanupChan <- "scraper"
			}()
		}

		if a.IsOn(app.Api) {
			runningServices++
			apiSvc := services.NewApiService(a.Logger)
			go func() {
				services.StartService(apiSvc, nil)
				cleanupChan <- apiSvc.Name()
			}()
			a.Logger.Info(
				"api is running",
				"apiUrl", apiSvc.ApiUrl(),
			)
		}

		if a.IsOn(app.Ipfs) {
			runningServices++
			ipfsSvc := services.NewIpfsService(a.Logger)
			go func() {
				services.StartService(ipfsSvc, nil)
				cleanupChan <- ipfsSvc.Name()
			}()
			a.Logger.Info(
				"ipfs daemon is running",
				"apiPort", ipfsSvc.ApiPort(),
				"apiMultiaddr", ipfsSvc.ApiMultiaddr(),
				"wasRunning", ipfsSvc.WasRunning(),
			)
		}

		if a.IsOn(app.Monitor) {
			runningServices++
			monSvc := services.NewMonitorService(a.Logger)
			go func() {
				services.StartService(monSvc, nil)
				cleanupChan <- monSvc.Name()
			}()
			a.Logger.Info(
				"monitors are running",
			)
		}

		<-stopChan
		// a.Logger.Info("received Control+C, waiting for all services to clean up...")

		for i := 0; i < runningServices; i++ {
			<-cleanupChan
			// a.Logger.Info("service cleanup complete", "service", serviceName)
		}

		// a.Logger.Info("all services cleaned up, exiting.")
	}
}
