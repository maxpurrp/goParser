package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

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
	pages   int = 10
	headers     = read_head_value()
	wg      sync.WaitGroup
)

type Car struct {
	PlateNumber string `json:"PlateNumber"`
	BigPhotoURL string `json:"BigPhotoURL"`
	PlatePhoto  string `json:"PlatePhoto"`
	Model       string `json:"Model"`
}

func read_head_value() []string {
	var head_values []string
	file, err := os.Open("./head_value.txt")
	check_err(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		head_values = append(head_values, line)
	}
	return head_values
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func _get_car_model(doc *goquery.Document) string {
	html := doc.Find(".col-md-6.col-sm-7").Find(".panel-body").Find(".text-center.margin-bottom-10")
	carModel := strings.TrimSpace(html.Text())
	return carModel
}

func _get_car_fullPhoto(doc *goquery.Document) string {
	html := doc.Find(".col-md-6.col-sm-7").Find(".panel-body a")
	fullPhoto, _ := html.Attr("href")
	return fullPhoto
}

func _get_car_platePhoto(doc *goquery.Document) string {
	platePhoto, _ := doc.Find(".col-md-6.col-sm-7").Find(".panel-body").Find(".img-responsive.center-block.margin-bottom-20").Attr("src")
	return platePhoto
}

func _get_car_number(doc *goquery.Document) string {
	html := doc.Find(".breadcrumbs").Find(".col-xs-12 h1")
	carNumber := strings.TrimSpace(html.Text())
	return carNumber
}

func send_req(url string) *http.Response {
	var res *http.Response
	for i := 0; i < len(headers); i++ {
		req, err := http.NewRequest("GET", url, nil)
		check_err(err)
		req.Header.Set("User-Agent", headers[i])
		req.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml;q=0.8, image/webp, */*;q=0.9")
		res, err = client.Do(req)
		check_err(err)
		if res.StatusCode == http.StatusOK {
			return res
		}
	}
	if res.StatusCode != http.StatusOK {
		res = send_req(url)
	}
	return res
}

func getBody(country string) {
	var response *http.Response
	for i := 0; i < pages; i++ {
		var URL string
		if i == 0 {
			URL = fmt.Sprintf("https://platesmania.com/%s/gallery", country)
		} else {
			URL = fmt.Sprintf("https://platesmania.com/%s/gallery-%d", country, i)
		}
		response = send_req(URL)
		defer response.Body.Close()
		links := getLinks(response, country)
		go parseCard(links, country, i+1)
	}
}

func getLinks(body *http.Response, country string) []string {
	//PARSE BODY AND GET LINKS FOR CARS
	var links []string
	doc, err := goquery.NewDocumentFromReader(body.Body)
	check_err(err)
	selector := ".row.blog-page"
	doc.Find(selector).Find(".panel-body").Find(".row").Find(".col-xs-offset-3 a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		links = append(links, link)
	})
	return links
}

func parseCard(links []string, country string, page int) {
	for i := 0; i < len(links); i++ {
		wg.Add(1)
		defer wg.Done()
		URL := fmt.Sprintf("https://platesmania.com%s", links[i])
		response := send_req(URL)
		defer response.Body.Close()
		doc, err := goquery.NewDocumentFromReader(response.Body)
		check_err(err)
		car := Car{
			Model:       _get_car_model(doc),
			PlateNumber: _get_car_number(doc),
			BigPhotoURL: "https://platesmania.com" + _get_car_fullPhoto(doc),
			PlatePhoto:  _get_car_platePhoto(doc),
		}
		path := country + "/" + fmt.Sprint(page)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		jsonData, err := json.MarshalIndent(car, "", " ")
		check_err(err)
		file, err := os.Create(path + "/" + fmt.Sprint(i+1) + ".json")
		check_err(err)
		defer file.Close()
		_, err = file.Write(jsonData)
		check_err(err)
	}
}

func main() {
	countries := []string{"us", "ru", "ua", "uk", "fr"}
	for i := 0; i < len(countries); i++ {
		wg.Add(1)
		go func(country string) {
			defer wg.Done()
			getBody(country)
		}(countries[i])
	}
	wg.Wait()
	fmt.Println("Done")
}
