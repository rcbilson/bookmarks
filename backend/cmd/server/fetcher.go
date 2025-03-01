package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type BookmarkData struct {
	Title string
	Icon  []byte
}

type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
	FetchBookmark(ctx context.Context, url string) (BookmarkData, error)
}

type FetcherImpl struct {
}

func NewFetcher() (Fetcher, error) {
	return &FetcherImpl{}, nil
}

func (*FetcherImpl) Fetch(ctx context.Context, url string) ([]byte, error) {
	var httpClient http.Client

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// spoof user agent to work around bot detection
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36"}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Println("Headers:")
		for k, v := range res.Header {
			log.Println("    ", k, ":", v)
		}
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func findChild(n *html.Node, dataAtom atom.Atom) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.DataAtom == dataAtom {
			return c
		}
	}
	return nil
}

func (fetcher *FetcherImpl) FetchBookmark(ctx context.Context, url string) (bookmark BookmarkData, err error) {
	page, err := fetcher.Fetch(ctx, url)
	if err != nil {
		return bookmark, fmt.Errorf("Error retrieving site: %v", err)
	}
	// parse the html in the page and extract the title
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return
	}

	htmlNode := findChild(doc, atom.Html)
	if htmlNode == nil {
		return
	}
	headNode := findChild(htmlNode, atom.Head)
	if headNode == nil {
		return
	}
	for n := headNode.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode && n.DataAtom == atom.Title {
			bookmark.Title = n.FirstChild.Data
			break
		}
	}
	return
}
