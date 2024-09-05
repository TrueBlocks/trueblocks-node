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
		wg.Add(3)
		go a.serve(&wg)
		go a.scrape(&wg)
		go a.monitor(&wg)
		wg.Wait()
	}
}

var screenMutex sync.Mutex
