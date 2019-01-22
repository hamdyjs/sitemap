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

	for i, linkStr := range links {
		childLink, err := url.Parse(linkStr.Href)
		if err != nil {
			continue
		}
		referenceLink := pageLink.ResolveReference(childLink)
		if referenceLink.Host != pageLink.Host {
			links = append(links[:i], links[i+1:]...)
			continue
		}
		childLinks, err := getLinksFromURL(referenceLink)
		if err != nil {
			continue
		}
		links = append(links, childLinks...)
	}

	oldLinks := make([]link.Link, len(links))
	copy(oldLinks, links)
	links = make([]link.Link, 0)
	linksDone := make(map[string]bool)
	for _, link := range oldLinks {
		fmt.Println(link.Href, linksDone[link.Href])
		if !linksDone[link.Href] {
			links = append(links, link)
		}
		linksDone[link.Href] = true
	}

	return links, nil
}
