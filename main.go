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

var linksFetched = make(map[string]bool)

func getLinksFromURL(pageLink *url.URL) ([]link.Link, error) {
	if linksFetched[pageLink.String()] {
		return nil, nil
	}
	linksFetched[pageLink.String()] = true
	res, err := http.Get(pageLink.String())
	if err != nil {
		return nil, err
	}

	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	externalLinks := make(map[string]bool)

	for _, link := range links {
		childLink, err := url.Parse(link.Href)
		if err != nil {
			continue
		}
		referenceLink := pageLink.ResolveReference(childLink)
		if referenceLink.Host != pageLink.Host {
			externalLinks[link.Href] = true
			continue
		}
		childLinks, err := getLinksFromURL(referenceLink)
		if err != nil {
			continue
		}
		links = append(links, childLinks...)
	}

	// Filter duplicate and external links
	linksDone := make(map[string]bool)
	finalLinks := links[:0]
	for _, link := range links {
		if !linksDone[link.Href] && !externalLinks[link.Href] {
			finalLinks = append(finalLinks, link)
			linksDone[link.Href] = true
		}
	}

	return finalLinks, nil
}
