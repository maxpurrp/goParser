package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	carLinks = make(map[string][]string)
)

type Car struct {
	NumberPlate   string
	BigPhotoURL   string
	SmallPhotoURL string
	Model         string
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getBody(country string) {
	for i := 0; i < 10; i++ {
		var URL string
		if i == 0 {
			URL = fmt.Sprintf("https://platesmania.com/%s/gallery", country)
		} else {
			URL = fmt.Sprintf("https://platesmania.com/%s/gallery-%d", country, i)
		}
		req, err := http.NewRequest("GET", URL, nil)
		check_err(err)
		req.Header.Set("User-Agent", "go:getter")
		res, err := client.Do(req)
		check_err(err)
		defer res.Body.Close()
		links := ParseBody(res, country)
		carLinks[country] = append(carLinks[country], links...)
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

func parseCard(link string) {
	URL := fmt.Sprintf("https://platesmania.com%s", link)
	req, err := http.NewRequest("GET", URL, nil)
	check_err(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	check_err(err)
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
		Model:         carModel,
		NumberPlate:   carNumber,
		BigPhotoURL:   "https://platesmania.com" + fullPhoto,
		SmallPhotoURL: platePhoto,
	}
	fmt.Println(car.NumberPlate)
	fmt.Println(car.BigPhotoURL)
	fmt.Println(car.SmallPhotoURL)
	fmt.Println(car.Model)

}

func main() {
	countries := []string{"ru", "ua", "uz", "cn", "us"}
	for _, country := range countries {
		getBody(country)
	}
	for _, links := range carLinks {
		for _, link := range links {
			parseCard(link)
		}
	}

}
