package main

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func main() {
	establishConfig()

	if err := startApiServer(); err != nil {
		logger.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go scraper(&wg)
	go monitor(&wg)
	wg.Wait()
}

var screenMutex sync.Mutex
