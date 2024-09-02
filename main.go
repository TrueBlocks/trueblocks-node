package main

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func main() {
	if config, err := establishConfig(); err != nil {
		logger.Error(err)
	} else {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go scraper(config, &wg)
		// go monitor(config, &wg)
		wg.Wait()
	}
}

var screenMutex sync.Mutex
