package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type Card struct {
	Title   string `json:"title"`
	Price   string `json:"price"`
	Set     string `json:"set"`
	InStock bool   `json:"instock"`
}

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

	priceFromPage := doc.Find(".price").Text()
	stockFromPage := doc.Find(".product-form__add-button").Text()
	titleFromPage := doc.Find(".product-meta__title").Text()

	stockStatus := false

	if stockFromPage != "Add to cart" {
		stockStatus = false
	} else {
		stockStatus = true
	}

	cardName := strings.Split(titleFromPage, " :: ")
	testPrice := strings.Split(priceFromPage, "$")

	card := Card{
		Title: cardName[0], Price: testPrice[1], Set: cardName[1], InStock: stockStatus,
	}

	json.NewEncoder(w).Encode((card))

	fmt.Println(cardName[0], testPrice[1], cardName[1], stockFromPage)

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
