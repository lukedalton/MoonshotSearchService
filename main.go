package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	cardName := "liliana-of-the-veil-uma"

	// Request the HTML page.
	res, err := http.Get("https://moonshotgamestore.com/products/" + cardName)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	price, e := regexp.Compile(`\$\d+\.\d+`)

	CheckError(e)

	title := doc.Find(".price").Text()
	stock := doc.Find(".product-form__add-button").Text()

	if stock != "Add to cart" {
		stock = "Out of Stock"
	} else {
		stock = "In stock"
	}

	trimmed := price.Find([]byte(title))

	fmt.Println("Liliana of the Veil:", string(trimmed), stock)

}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func main() {
	ExampleScrape()
}
