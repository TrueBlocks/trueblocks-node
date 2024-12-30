package main

import (
	"github.com/TrueBlocks/trueblocks-node/v4/app"
	"github.com/TrueBlocks/trueblocks-sdk/v4/services"
)

func main() {
	a := app.NewApp()
	if err := a.EstablishConfig(); err != nil {
		a.Fatal(err)
	}

	if cont, activeServices, err := a.ParseArgs(); !cont {
		return
	} else if err != nil {
		a.Fatal(err)
	} else {
		a.Logger.Info("Starting TrueBlocks Node", "services", len(activeServices))

		serviceManager := services.NewServiceManager(activeServices, a.Logger)

		for _, svc := range activeServices {
			if controlSvc, ok := svc.(*services.ControlService); ok {
				controlSvc.AttachServiceManager(serviceManager)
			}
		}

		if err := serviceManager.StartAllServices(); err != nil {
			a.Fatal(err)
		}

		serviceManager.HandleSignals()

		select {}
	}
}
