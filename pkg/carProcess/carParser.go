package carProcess

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Car struct {
	PlateNumber string `json:"PlateNumber"`
	BigPhotoURL string `json:"BigPhotoURL"`
	PlatePhoto  string `json:"PlatePhoto"`
	Model       string `json:"Model"`
}

func (c *Car) getCarChars(doc goquery.Document) {
	selector := ".col-md-6.col-sm-7 .panel-body "
	c.Model = c.getDataFromDoc(&doc, selector+".text-center.margin-bottom-10", "")
	c.BigPhotoURL = c.getDataFromDoc(&doc, selector+"a", "href")
	c.PlatePhoto = c.getDataFromDoc(&doc, selector+".img-responsive.center-block.margin-bottom-20", "src")
	c.PlateNumber = c.getDataFromDoc(&doc, ".breadcrumbs .col-xs-12 h1", "")
}

func (c *Car) getDataFromDoc(doc *goquery.Document, selector string, attr string) string {
	result := doc.Find(selector)
	if attr != "" {
		value, _ := result.Attr(attr)
		return value
	}
	return strings.TrimSpace(result.Text())
}

func (c *Car) saveToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file : %v", err)
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("failed to Marshal json : %v", err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to file : %v", err)
	}
	return nil
}

func ManageCarData(response *http.Response, country string, page int) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Printf("failed read document : %v", err)
	}

	car := Car{}
	car.getCarChars(*doc)
	//creating a folder in the application root
	dir := path.Join("data", country, fmt.Sprint(page))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	//Some numbers have a '/' sign, which conflicts with the creation of folders
	name := strings.Replace(car.PlateNumber, "/", "-", -1)
	filePath := path.Join(dir, fmt.Sprint(name)+".json")
	//save to file
	err = car.saveToFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
}
