package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func scraper(wg *sync.WaitGroup) {
	defer wg.Done()
	opts := sdk.InitOptions{}
	if _, meta, err := opts.Init(); err != nil {
		logger.Error("Error:", err)
		return
	} else {
		fmt.Println("Meta:", meta.String())
	}

	for {
		fmt.Println(colors.Green, "Scraper is running", colors.Off)
		opts := sdk.ScrapeOptions{}
		if blocks, meta, err := opts.ScrapeRunCount(1); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Blocks:", blocks)
			fmt.Println("Meta:", meta)
		}
		fmt.Println(os.Getpid())
		time.Sleep(time.Second * 3)
	}
}
