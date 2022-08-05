package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type MetaData struct {
	MetaTitle       string `json:"title"`
	MetaDescription string `json:"description"`
	MetaKeywords    string `json:"keywords"`
	MetaImage       string `json:"image"`
}

type ReqBody struct {
	WebsiteUrl string `json:"website_url"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get-metadata", MetaHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":9000", r))
}

func MetaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t ReqBody
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		panic(err)
	}

	res := metaScrape(t.WebsiteUrl)

	json.NewEncoder(w).Encode(res)
}

func metaScrape(url string) MetaData {
	// Request the HTML page.
	res, error := http.Get(url)

	if error != nil {
		log.Fatal(error)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var mDescription string
	var mKeywords string
	var mImage string
	pageTitle := doc.Find("title").Contents().Text()

	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" {
			mDescription = item.AttrOr("content", "")
		}
	})
	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "keywords" {
			mKeywords = item.AttrOr("content", "")
		}
	})
	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("property", "") == "og:image" {
			mImage = item.AttrOr("content", "")
		}
	})
	if mImage == "" {
		doc.Find("img").Each(func(index int, item *goquery.Selection) {
			if index == 0 {
				mImage = item.AttrOr("src", "")
			}
		})
	}

	finalMeta := MetaData{pageTitle, mDescription, mKeywords, mImage}

	return finalMeta

}
