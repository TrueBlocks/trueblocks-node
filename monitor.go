package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

func monitor(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Println(colors.Green, "Monitor is running", colors.Off)
		time.Sleep(time.Second * 3)
	}
}
