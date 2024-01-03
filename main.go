package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	tr *http.Transport = &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
		},
	}
	client *http.Client = &http.Client{
		Transport: tr,
	}
	headers = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.345",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.23",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.21",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/91.0.864.37",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/91.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Android 11; Mobile; rv:91.0) Gecko/91.0 Firefox/91.0"}
)

type Car struct {
	PlateNumber string `json:"PlateNumber"`
	BigPhotoURL string `json:"BigPhotoURL"`
	PlatePhoto  string `json:"PlatePhoto"`
	Model       string `json:"Model"`
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func send_req(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	check_err(err)
	req.Header.Set("User-Agent", headers[0])
	req.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, */*;q=0.8")
	res, err := client.Do(req)
	check_err(err)
	if res.StatusCode != http.StatusOK {
		for i := 1; i < len(headers)-1; i++ {
			req, err := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", headers[i])
			req.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml;q=0.8, image/webp, */*;q=0.9")
			res, err := client.Do(req)
			check_err(err)
			if res.StatusCode == http.StatusOK {
				return res
			} else {
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
	fmt.Println(res.Status)
	return res
}

func getBody(country string) {
	var res *http.Response
	for i := 0; i < 10; i++ {
		var URL string
		if i == 0 {
			URL = fmt.Sprintf("https://platesmania.com/%s/gallery", country)
		} else {
			URL = fmt.Sprintf("https://platesmania.com/%s/gallery-%d", country, i)
		}
		res = send_req(URL)
		defer res.Body.Close()
		links := ParseBody(res, country)
		go parseCard(links, country, i+1)
	}
}

func ParseBody(body *http.Response, country string) []string {
	var links []string
	doc, err := goquery.NewDocumentFromReader(body.Body)
	check_err(err)
	selector := ".row.blog-page"
	doc.Find(selector).Find(".panel-body").Find(".row").Find(".col-xs-offset-3 a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			links = append(links, link)
		}
	})
	return links
}

func parseCard(links []string, country string, page int) {
	for i := 0; i < len(links); i++ {
		URL := fmt.Sprintf("https://platesmania.com%s", links[i])
		res := send_req(URL)
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		check_err(err)

		html := doc.Find(".col-md-6.col-sm-7").Find(".panel-body a")
		fullPhoto, _ := html.Attr("href")

		html = doc.Find(".col-md-6.col-sm-7").Find(".panel-body").Find(".text-center.margin-bottom-10")
		carModel := strings.TrimSpace(html.Text())

		platePhoto, _ := doc.Find(".col-md-6.col-sm-7").Find(".panel-body").Find(".img-responsive.center-block.margin-bottom-20").Attr("src")

		html = doc.Find(".breadcrumbs").Find(".col-xs-12 h1")
		carNumber := strings.TrimSpace(html.Text())
		car := Car{
			Model:       carModel,
			PlateNumber: carNumber,
			BigPhotoURL: "https://platesmania.com" + fullPhoto,
			PlatePhoto:  platePhoto,
		}
		path := country + "/" + fmt.Sprint(page)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		jsonData, err := json.MarshalIndent(car, "", " ")
		check_err(err)
		file, err := os.Create(path + "/" + fmt.Sprint(i+1) + ".json")
		defer file.Close()
		check_err(err)
		_, err = file.Write(jsonData)
		check_err(err)
	}

}

func main() {
	countries := []string{"us", "ru", "ua", "uk", "fr"}
	for i := 0; i < len(countries); i++ {
		getBody(countries[i])
	}
	for i := 0; i != 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("still working")
	}
}
