package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/url"
	"os"
	"time"
)

type Product struct {
	title        string
	currentPrice string
	oldPrice     string
	url          string
	discountType string
}

type Products []Product

var searchList = []string{
	"notebook",
	"gaming mouse",
}

const domain = "https://www.amazon.com.tr"
const pageLimit = 5

func startCrawl() {

	var products []Product

	c := colly.NewCollector(colly.Async(true))

	// Random delay between request
	c.Limit(&colly.LimitRule{
		RandomDelay: 3 * time.Second,
		Parallelism: 6,
	})

	// Set random user agent
	extensions.RandomUserAgent(c)

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("URL: ", request.URL)
		fmt.Println("User Agent: ", request.Headers.Get("User-Agent"))
	})

	c.OnHTML("div.s-main-slot.s-result-list.s-search-results.sg-row", func(element *colly.HTMLElement) {
		element.ForEach("div.a-section.a-spacing-medium", func(_ int, element *colly.HTMLElement) {
			temp := Product{}
			temp.title = element.ChildText("span .a-size-base-plus.a-color-base.a-text-normal")
			temp.currentPrice = element.ChildText("span[data-a-color='base'] .a-offscreen")
			temp.oldPrice = element.ChildText("span[data-a-color='secondary'] .a-offscreen")
			temp.discountType = element.ChildText("span[data-a-badge-type='deal']")
			var slug, _ = element.DOM.Find(".a-link-normal.a-text-normal").Attr("href")

			if slug != "" {
				temp.url = fmt.Sprintf("%s%s", domain, slug)
			}

			if temp.title == "" || temp.currentPrice == "" || temp.discountType == "" {
				return
			}

			products = append(products, temp)
		})
	})

	for k := 0; k < len(searchList); k++ {
		fmt.Println(searchList[k])
		for i := 1; i <= pageLimit; i++ {
			visitUrl := fmt.Sprintf("%s/s?k=%s&page=%d", domain, url.QueryEscape(searchList[k]), i)
			c.Visit(visitUrl)
		}
	}

	c.Wait()

	exportCSV(products)
}

func exportCSV(products Products) {
	file, err := os.Create("results.csv")
	defer file.Close()
	if err != nil {
		log.Fatalln("Failed to open file!", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, product := range products {
		row := []string{product.title, product.url, product.oldPrice, product.currentPrice}
		data = append(data, row)
	}
	w.WriteAll(data)
}

func main() {
	startCrawl()
}
