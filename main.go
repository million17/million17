package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type quote struct {
	Text   string `json:"quote"`
	Author string `json:"author"`
}

const quoteAPI = "https://dummyjson.com/quotes/random"

func fetchQuote(client *http.Client) (quote, error) {
	request, err := http.NewRequest(http.MethodGet, quoteAPI, nil)
	if err != nil {
		return quote{}, fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "million17-daily-quote/1.0")

	response, err := client.Do(request)
	if err != nil {
		return quote{}, fmt.Errorf("request quote API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return quote{}, fmt.Errorf("quote API returned %s: %s", response.Status, strings.TrimSpace(string(body)))
	}

	var selected quote
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&selected); err != nil {
		return quote{}, fmt.Errorf("decode quote API response: %w", err)
	}
	selected.Text = strings.TrimSpace(selected.Text)
	selected.Author = strings.TrimSpace(selected.Author)
	if selected.Text == "" || selected.Author == "" {
		return quote{}, fmt.Errorf("quote API returned an empty quote or author")
	}

	return selected, nil
}

func main() {
	client := &http.Client{Timeout: 15 * time.Second}
	selected, err := fetchQuote(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	content := fmt.Sprintf("> %s\n\n**%s**\n", selected.Text, selected.Author)
	if err := os.WriteFile("README.md", []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "update README.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Selected quote by %s\n", selected.Author)
}
