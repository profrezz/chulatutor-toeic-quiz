package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RowData struct {
	Eng      string `json:"eng"`
	Meaning  string `json:"meaning"`
	Sentence string `json:"sentence"`
}

func main() {
	url := "https://www.chulatutor.com/blog/%E0%B8%84%E0%B8%B3%E0%B8%A8%E0%B8%B1%E0%B8%9E%E0%B8%97%E0%B9%8C-toeic/"

	// Make HTTP request
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch data: status code %d", res.StatusCode)
	}

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	doc.Find("table").Each(func(tableIndex int, table *goquery.Selection) {
		var rows []RowData

		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			cells := row.Find("td")

			if cells.Length() >= 3 {
				data := RowData{
					Eng:      strings.TrimSpace(cells.Eq(0).Text()),
					Meaning:  strings.TrimSpace(cells.Eq(1).Text()),
					Sentence: strings.TrimSpace(cells.Eq(2).Text()),
				}
				rows = append(rows, data)
			}
		})

		if len(rows) > 0 {
			filename := fmt.Sprintf("table_%02d.json", tableIndex+1)
			file, err := os.Create(filename)
			if err != nil {
				log.Printf("Error creating file %s: %v", filename, err)
				return
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(rows); err != nil {
				log.Printf("Error writing JSON to file: %v", err)
			} else {
				fmt.Printf("Saved %d rows to %s\n", len(rows), filename)
			}
		}
	})
}
