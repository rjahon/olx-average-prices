package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type DestURL struct {
	City            string
	NORoomsFrom     int
	NORoomsTo       int
	TotalAreaFrom   int
	TotalAreaTo     int
	FloorFrom       int
	FloorTo         int
	TotalFloorsFrom int
	TotalFloorsTo   int
}

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

	// var dest destURL
	destKZ := DestURL{
		City:            "alma-ata",
		NORoomsFrom:     2,
		NORoomsTo:       2,
		TotalAreaFrom:   40,
		TotalAreaTo:     80,
		FloorFrom:       2,
		FloorTo:         3,
		TotalFloorsFrom: 4,
		TotalFloorsTo:   7,
	}

	urlKZ := fmt.Sprintf(
		"https://www.olx.kz/d/nedvizhimost/kvartiry/prodazha/%s/?search[filter_float_number_of_rooms:from]=%d&search[filter_float_number_of_rooms:to]=%d&search[filter_float_total_area:from]=%d&search[filter_float_total_area:to]=%d&search[filter_float_floor:from]=%d&search[filter_float_floor:to]=%d&search[filter_float_total_floors:from]=%d&search[filter_float_total_floors:to]=%d",
		destKZ.City,
		destKZ.NORoomsFrom,
		destKZ.NORoomsTo,
		destKZ.TotalAreaFrom,
		destKZ.TotalAreaTo,
		destKZ.FloorFrom,
		destKZ.FloorTo,
		destKZ.TotalFloorsFrom,
		destKZ.TotalFloorsTo,
	)

	c.Visit(urlKZ)

	fmt.Printf("Average price: %d", avgPrice(prices))
	outJSON(prices)
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

func outJSON(prices []int) {
	content, err := json.Marshal(prices)
	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile("data.json", content, 0644)
}
