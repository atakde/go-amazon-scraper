package main

import (
	"amazon-scraper-go/helper"
	"amazon-scraper-go/model"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"sort"
	"time"
)

var searchList = []string{
	"notebook",
	"gaming mouse",
}

const domain = "https://www.amazon.com"
const pageLimit = 5

func startCrawl() {

	var products []model.Product

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
			temp := model.Product{}
			temp.Title = element.ChildText("span .a-size-base-plus.a-color-base.a-text-normal")
			temp.CurrentPrice = helper.FormatPrice(element.ChildText("span[data-a-color='base'] .a-offscreen"))
			temp.OldPrice = helper.FormatPrice(element.ChildText("span[data-a-color='secondary'] .a-offscreen"))
			temp.DiscountType = element.ChildText("span[data-a-badge-type='deal']")
			var slug, _ = element.DOM.Find(".a-link-normal.a-text-normal").Attr("href")

			if slug != "" {
				temp.Url = fmt.Sprintf("%s%s", domain, slug)
			}

			if temp.Title == "" {
				return
			}

			if temp.DiscountType != "" {
				products = append(products, temp)
			}
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

	//helper.ExportCSV(products)

	// sort
	sort.Slice(products, func(p, q int) bool {
		return products[p].CurrentPrice > products[q].CurrentPrice
	})

	currentTime := time.Now().Local()
	subject := "Deal Scraper Report | " + currentTime.Format("2006-01-02 15:04:05")
	to := os.Getenv("EMAIL_TO")
	from := os.Getenv("EMAIL_FROM")
	toName := "Atakan"
	fromName := "Deal Scraper"
	r := helper.NewMail(to, toName, from, fromName, subject)

	r.Send("static/templates/mail.html", products)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	startCrawl()
}
