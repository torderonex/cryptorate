package rubparser

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const parsingURL = "https://cbr.ru/currency_base/daily/"

func Parse() float64 {

	resp, err := http.Get(parsingURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0
	}

	defer resp.Body.Close()
	var res string
	var isNext bool
	var temp func(*html.Node)
	temp = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "td" {
			if isNext {
				res = node.FirstChild.Data
				isNext = false
				return
			}
			if node.FirstChild != nil && (node.FirstChild.Data == "US Dollar" || node.FirstChild.Data == "Доллар США") {
				isNext = true
				fmt.Println(node.NextSibling.Data)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			temp(c)
		}
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return 0
	}
	temp(doc)
	res = strings.ReplaceAll(res, ",", ".")
	course, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return 0
	}

	return course
}
