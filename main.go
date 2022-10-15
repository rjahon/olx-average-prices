package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains(
			"olx.com",

			"olx.kz",
			"www.olx.kz",
		),
	)

	var prices []int

	c.OnHTML("div.css-1venxj6", func(e *colly.HTMLElement) {
		price, err := cleanStr(e.ChildText("p[data-testid=ad-price]"))
		if err == nil {
			prices = append(prices, price)
		}
	})

	c.OnHTML("a[data-testid=pagination-forward]", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n\n", r.URL.String())
	})

	// kz
	c.Visit("https://www.olx.kz/d/nedvizhimost/kvartiry/prodazha/alma-ata/?search%5Bfilter_float_number_of_rooms:from%5D=2&search%5Bfilter_float_number_of_rooms:to%5D=2&search%5Bfilter_float_floor:to%5D=2")

	fmt.Printf("Average price: %d", avgPrice(prices))
}

func cleanStr(s string) (int, error) {
	for k, v := range s {
		if v == 'т' {
			s = s[:k]
		}
	}
	s = strings.ReplaceAll(s, " ", "")
	price, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func avgPrice(prices []int) int {
	count := 0
	total_price := 0
	for _, price := range prices {
		total_price += price
		count++
	}

	if count > 0 {
		return total_price / count
	} else {
		return 0
	}
}
