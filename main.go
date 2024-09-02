package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("v3.0.3-2024-09-01-23-18-10")
		return
	}

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
