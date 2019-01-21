package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hamdyjs/link"
)

func main() {
	urlString := flag.String("url", "http://gophercises.com", "The url to build the sitemap of")
	flag.Parse()

	url, err := url.Parse(*urlString)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	siteLinks, err := getLinksFromURL(url.Path)
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
