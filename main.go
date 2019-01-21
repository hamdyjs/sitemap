package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/hamdyjs/link"
)

func main() {
	url := flag.String("url", "http://gophercises.com", "The url to build the sitemap of")
	flag.Parse()

	siteLinks, err := getLinksFromURL(*url)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Printf("%+v", siteLinks)
}

func getLinksFromURL(url string) ([]link.Link, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	return links, nil
}
