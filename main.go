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

	res, err := http.Get(*url)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	siteLinks, err := link.Parse(res.Body)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	res.Body.Close()

	fmt.Printf("%+v", siteLinks)
}
