package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

func CardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	res, err := http.Get("https://moonshotgamestore.com/products/" + vars["cardName"])
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	regExpForprice, e := regexp.Compile(`\$\d+\.\d+`)
	regExpForCardName, e := regexp.Compile(`(.)*`)

	CheckError(e)

	priceFromPage := doc.Find(".price").Text()
	stockFromPage := doc.Find(".product-form__add-button").Text()
	titleFromPage := doc.Find(".product-meta__title").Text()

	if stockFromPage != "Add to cart" {
		stockFromPage = "Out of Stock"
	} else {
		stockFromPage = "In stock"
	}

	trimmedPrice := regExpForprice.Find([]byte(priceFromPage))
	trimmedTitle := regExpForCardName.Find([]byte(titleFromPage))

	fmt.Println(string(trimmedTitle), string(trimmedPrice), stockFromPage)

}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/card/{cardName}", CardHandler)
	http.Handle("/", router)

	http.ListenAndServe(":8000", router)
}
