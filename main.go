package main

import (
	"sync"
)

func main() {
	if a, err := NewApp(); err != nil {
		a.Logger.Error(err.Error())
		return
	} else {
		if cont := a.ParseArgs(); !cont {
			return
		}
		wg := sync.WaitGroup{}
		wg.Add(2)
		go a.scraper(&wg)
		// go monitor(config, &wg)
		wg.Wait()
	}
}

var screenMutex sync.Mutex
