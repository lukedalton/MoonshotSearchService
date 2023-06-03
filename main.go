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

	// Currently Moonshot's URL plus /products/{Card Name} with "-"" replacing spaces in the card name, and three letter set identifer at the end.
	// If it is a foil, card-name-foil-set identifier.
	res, err := http.Get("https://moonshotgamestore.com/products/" + vars["cardName"])
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	if res.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		priceFromPage := doc.Find(".price").Text()
		stockFromPage := doc.Find(".product-form__add-button").Text()
		cardNameFromPage := doc.Find(".product-meta__title").Text()

		stockStatus := stockFromPage == "Add to cart"

		// Split on space colon colon space, so that first string in array is actual card name, and the second is set three letter abbreviation.
		cardName := strings.Split(cardNameFromPage, " :: ")
		// Used a split here, because it was the only think I can get to work correctly. The first item in array is useless, the second contains price.
		price := strings.Split(priceFromPage, "$")

		card := Card{
			Title: cardName[0], Price: price[1], Set: cardName[1], InStock: stockStatus,
		}

		json.NewEncoder(w).Encode((card))
	}

	if res.StatusCode == 404 {
		w.WriteHeader(res.StatusCode)
	}
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
