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
		screenMutex.Lock()
		fmt.Println("\r"+colors.Green+"Monitor is running", colors.Off)
		screenMutex.Unlock()
		time.Sleep(time.Millisecond * 3000)
	}
}
