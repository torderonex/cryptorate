package cryptocurrencyparser

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func Parse(url string) string {
	if !strings.HasPrefix(url, "https://coinmarketcap.com/currencies/") {
		url = "https://coinmarketcap.com/currencies/" + url
	}
	resp, _ := http.Get(url)
	body, _ := html.Parse(resp.Body)
	defer resp.Body.Close()
	var temp func(*html.Node)
	var res string
	temp = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "span" {
			for _, attr := range node.Parent.Attr {
				if attr.Key == "class" && attr.Val == "priceValue " {
					res = node.FirstChild.Data
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			temp(c)
		}
	}
	temp(body)
	return res
}
