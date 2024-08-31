package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func scraper(wg *sync.WaitGroup) {
	defer wg.Done()

	opts := sdk.InitOptions{}
	if _, _, err := opts.InitAll(); err != nil { // blooms only, if that fails
		if _, _, err := opts.InitAll(); err != nil { // try --all
			logger.Error(err)
			return
		}
	}

	cwd, _ := os.Getwd()
	dataFilename := filepath.Join(cwd, "meta.json")

	for {
		screenMutex.Lock()
		fmt.Print(colors.Green, "Scraper is running...", colors.Off)
		quit := false
		go func() {
			for {
				if quit {
					break
				}
				time.Sleep(time.Millisecond * 1000)
				fmt.Print(".")
			}
		}()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go scrapeOnce(dataFilename, &wg)
		wg.Wait()
		quit = true
		fmt.Println(colors.Green, "Done.", colors.Off)
		time.Sleep(time.Millisecond * 1000)
		fmt.Print("\r \r")
		screenMutex.Unlock()
		time.Sleep(time.Millisecond * 4000)
	}
}

func scrapeOnce(dataFilename string, wg *sync.WaitGroup) {
	defer wg.Done()
	opts := sdk.ScrapeOptions{
		BlockCnt: 500,
		Globals: sdk.Globals{
			Chain: "mainnet",
		},
	}
	w := logger.GetLoggerWriter()
	logger.SetLoggerWriter(io.Discard)
	if _, meta, err := opts.ScrapeRunCount(1); err != nil {
		logger.Error(err)
	} else {
		fmt.Println(strings.ReplaceAll(strings.ReplaceAll(meta.String(), "\n", ""), " ", ""))
		file.StringToAsciiFile(dataFilename, meta.String())
	}
	logger.SetLoggerWriter(w)
}
