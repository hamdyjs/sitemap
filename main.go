package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hamdyjs/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type xmlURL struct {
	Location string `xml:"loc"`
}

type urlset struct {
	Links     []xmlURL `xml:"url"`
	Namespace string   `xml:"xmlns,attr"`
}

func main() {
	urlString := flag.String("url", "http://gophercises.com", "The url to build the sitemap of")
	flag.Parse()

	siteURL, err := url.Parse(*urlString)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	siteURLs, err := getLinksFromURL(siteURL)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	for _, link := range siteURLs {
		fmt.Println(link)
	}

	urlset := urlset{Namespace: xmlns}
	for _, url := range siteURLs {
		urlset.Links = append(urlset.Links, xmlURL{url.String()})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(urlset); err != nil {
		panic(err)
	}
}

var urlsFetched = make(map[string]bool)

func getLinksFromURL(pageURL *url.URL) ([]*url.URL, error) {
	if urlsFetched[pageURL.String()] {
		return nil, nil
	}
	urlsFetched[pageURL.String()] = true
	res, err := http.Get(pageURL.String())
	if err != nil {
		return nil, err
	}

	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	urls := make([]*url.URL, 0)
	for _, link := range links {
		if url, err := url.Parse(link.Href); err == nil {
			referenceURL := pageURL.ResolveReference(url)
			urls = append(urls, referenceURL)
		}
	}

	externalURLS := make(map[string]bool)

	for _, childURL := range urls {
		referenceURL := pageURL.ResolveReference(childURL)
		if referenceURL.Host != pageURL.Host {
			externalURLS[referenceURL.String()] = true
			continue
		}
		childURLs, err := getLinksFromURL(referenceURL)
		if err != nil {
			continue
		}
		urls = append(urls, childURLs...)
	}

	// Filter duplicate and external links
	urlsDone := make(map[string]bool)
	finalURLs := urls[:0]
	for _, url := range urls {
		if !urlsDone[url.String()] && !externalURLS[url.String()] {
			finalURLs = append(finalURLs, url)
			urlsDone[url.String()] = true
		}
	}

	return finalURLs, nil
}
