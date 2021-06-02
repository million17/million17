package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	Success struct {
		Total int `json:"total"`
	} `json:"success"`
	Contents struct {
		Quotes []struct {
			Quote      string   `json:"quote"`
			Length     string   `json:"length"`
			Author     string   `json:"author"`
			Tags       []string `json:"tags"`
			Category   string   `json:"category"`
			Language   string   `json:"language"`
			Date       string   `json:"date"`
			Permalink  string   `json:"permalink"`
			ID         string   `json:"id"`
			Background string   `json:"background"`
			Title      string   `json:"title"`
		} `json:"quotes"`
	} `json:"contents"`
	Baseurl   string `json:"baseurl"`
	Copyright struct {
		Year int    `json:"year"`
		URL  string `json:"url"`
	} `json:"copyright"`
}

func main() {
	res, err := http.Get("https://quotes.rest/qod?language=en&quot")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("Response status:", res.Status)
	fmt.Println("Response body:", res.Body)

	response := &Response{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		panic(err)
	}
	quote := response.Contents.Quotes[0]
	fmt.Println("Quote", quote.Quote)
	fmt.Println("Author", quote.Author)

	f, err := os.Create("README.md")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("_**%s**_\n\n%s", quote.Quote, quote.Author))

	if err != nil {
		panic(err)
	}
}