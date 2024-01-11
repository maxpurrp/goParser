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

var (
	mainLink string = "https://platesmania.com"
)

type Car struct {
	//funcstion in struct
	PlateNumber string `json:"PlateNumber"`
	BigPhotoURL string `json:"BigPhotoURL"`
	PlatePhoto  string `json:"PlatePhoto"`
	Model       string `json:"Model"`
}

func (c *Car) setCarChars(doc *goquery.Document) {
	c.Model = c.getCarModel(doc)
	c.PlatePhoto = c.getCarPlatePhoto(doc)
	c.BigPhotoURL = c.getCarFullPhoto(doc)
	c.PlateNumber = c.getCarNumber(doc)
}

func (c *Car) getCarModel(doc *goquery.Document) string {
	selector := ".col-md-6.col-sm-7 .panel-body .text-center.margin-bottom-10"
	return c.getDataFromDoc(doc, selector, "")
}

func (c *Car) getCarFullPhoto(doc *goquery.Document) string {
	selector := ".col-md-6.col-sm-7 .panel-body a"
	return mainLink + c.getDataFromDoc(doc, selector, "href")
}

func (c *Car) getCarPlatePhoto(doc *goquery.Document) string {
	selector := ".col-md-6.col-sm-7 .panel-body .img-responsive.center-block.margin-bottom-20"
	return c.getDataFromDoc(doc, selector, "src")
}

func (c *Car) getCarNumber(doc *goquery.Document) string {
	selector := ".breadcrumbs .col-xs-12 h1"
	return c.getDataFromDoc(doc, selector, "")
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
	car.setCarChars(doc)
	//creating a folder in the application root
	dir := path.Join("..", "data", country, fmt.Sprint(page))
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
