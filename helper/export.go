package helper

import (
	"encoding/csv"
	"log"
	"os"
	"amazon-scraper-go/model"
)

func ExportCSV(products model.Products) {
	file, err := os.Create("results.csv")
	defer file.Close()
	if err != nil {
		log.Fatalln("Failed to open file!", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, product := range products {
		row := []string{product.Title, product.Url, product.OldPrice, product.CurrentPrice}
		data = append(data, row)
	}
	w.WriteAll(data)
}
