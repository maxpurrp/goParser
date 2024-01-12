package main

import (
	"fmt"
	vars "parser/pkg"
	"parser/pkg/web"
)

func main() {
	countries := []string{"us", "ua", "ru", "fr", "uk"}
	for i, country := range countries {
		vars.WgMain.Add(1)

		go func(country string) {
			defer vars.WgMain.Done()
			web.GetBody(country)
		}(country)

		if (i+1)%2 == 0 {
			vars.WgMain.Wait()
		}
	}
	vars.WgMain.Wait()
	fmt.Printf("Done")
}
