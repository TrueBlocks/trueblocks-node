package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

func (a *App) monitor(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if !a.Busy {
			fmt.Println("\r"+colors.Green+"Monitor is running", colors.Off)
			time.Sleep(time.Millisecond * 3000)
		}
	}
}
