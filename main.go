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

	siteURL, err := url.Parse(*urlString)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	siteLinks, err := getLinksFromURL(siteURL)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	for _, link := range siteLinks {
		fmt.Println(link)
	}
}

func getLinksFromURL(pageLink *url.URL) ([]link.Link, error) {
	res, err := http.Get(pageLink.String())
	if err != nil {
		return nil, err
	}

	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	for _, linkStr := range links {
		childLink, err := url.Parse(linkStr.Href)
		if err != nil {
			continue
		}
		referenceLink := pageLink.ResolveReference(childLink)
		if referenceLink.Host != pageLink.Host {
			continue
		}
		childLinks, err := getLinksFromURL(referenceLink)
		if err != nil {
			continue
		}
		links = append(links, childLinks...)
	}

	return links, nil
}
