package main

import (
	"sync"
	"time"
)

// RunMonitor is a function that runs in a goroutine to monitor the app.
func (a *App) RunMonitor(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		a.Logger.Info("Monitor is running...")
		time.Sleep(time.Millisecond * 3000)
	}
}
