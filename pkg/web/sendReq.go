package web

import (
	"fmt"
	"math/rand"
	"net/http"
	vars "parser/pkg"
	"parser/pkg/carProcess"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func SendReq(url string) *http.Response {
	var res *http.Response
	for {
		req, err := http.NewRequest("GET", url, nil)
		checkErr(err)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Edge/97.0.1072.71 Yandex/22.1.2.155 Safari/537.36")
		req.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml;q=0.8, image/webp, */*;q=0.9")
		res, err = vars.Client.Do(req)
		checkErr(err)
		switch res.StatusCode {
		case 200:
			return res
		case 429:
			time.Sleep(2 * time.Second)
		}
		randTime := rand.Intn(251) + 250
		time.Sleep(time.Duration(randTime) * time.Millisecond)
	}
}

func GetBody(country string) {
	var wg sync.WaitGroup
	for i := 0; i < vars.Pages; i++ {
		wg.Add(1)

		var URL string
		if i == 0 {
			URL = fmt.Sprintf(vars.MainLink+"%s/gallery", country)
		} else {
			URL = fmt.Sprintf(vars.MainLink+"%s/gallery-%d", country, i)
		}

		go func(URL string, i int) {
			defer wg.Done()
			response := SendReq(URL)
			defer response.Body.Close()

			links := getLinks(response, country)
			for _, link := range links {
				response = SendReq(vars.MainLink + link)
				carProcess.ManageCarData(response, country, i+1)
			}
		}(URL, i)

		if (i+1)%2 == 0 {
			wg.Wait()
		}
	}
}

func getLinks(response *http.Response, country string) []string {
	var links []string
	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)
	selector := ".row.blog-page"
	doc.Find(selector).Find(".panel-body .row .col-xs-offset-3 a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		links = append(links, link)
	})
	return links
}
