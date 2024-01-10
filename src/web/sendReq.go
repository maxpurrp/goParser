package web

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	carProcess "parser/carProcess"

	"github.com/PuerkitoBio/goquery"
)

var (
	tr *http.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
		},
	}
	client *http.Client = &http.Client{
		Transport: tr,
	}
	pages        int    = 10
	mainLink     string = "https://platesmania.com/"
	wg           sync.WaitGroup
	maxCountGour int = 2
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
		res, err = client.Do(req)
		checkErr(err)
		if res.StatusCode == http.StatusOK {
			return res
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func GetBody(country string, wgMain *sync.WaitGroup, chMain chan struct{}) {
	defer wgMain.Done()

	ch := make(chan struct{}, maxCountGour)
	defer close(ch)

	for i := 0; i < pages; i++ {
		wg.Add(1)

		ch <- struct{}{}

		var URL string
		if i == 0 {
			URL = fmt.Sprintf(mainLink+"%s/gallery", country)
		} else {
			URL = fmt.Sprintf(mainLink+"%s/gallery-%d", country, i)
		}

		go func(URL string, i int) {
			defer wg.Done()
			response := SendReq(URL)
			defer response.Body.Close()

			links := getLinks(response, country)
			for _, link := range links {
				response = SendReq(mainLink + link)
				carProcess.ManageCarData(response, country, i+1)
			}
			<-ch
		}(URL, i)

	}
	<-chMain
	wg.Wait()
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
