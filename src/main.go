package main

import (
	"parser/web"
	"sync"
)

var (
	maxCountGour int = 2
)

func main() {
	countries := []string{"us", "ua", "ru", "fr", "uk"}
	var wg sync.WaitGroup
	ch := make(chan struct{}, maxCountGour)
	for i := 0; i < len(countries); i++ {
		ch <- struct{}{}
		wg.Add(1)
		go web.GetBody(countries[i], &wg, ch)
	}
	wg.Wait()
}
