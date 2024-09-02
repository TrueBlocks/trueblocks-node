package main

import "sync"

func main() {
	establishConfig()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go scraper(&wg)
	// go monitor(&wg)
	wg.Wait()
}

var screenMutex sync.Mutex
