package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

var urls = [...]string{
	"https://www.google.com",
	"https://www.seriouseats.com",
}

var titles = [...]string{
	"Google",
	"Serious Eats",
}

func TestFetch(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	fetcher, err := NewFetcher()
	assert.NilError(t, err)

	for _, url := range urls {
		bytes, err := fetcher.Fetch(context.Background(), url)
		if err != nil {
			t.Errorf("Failed to fetch %s", url)
		}

		// save files for other tests
		base := filepath.Base(url)
		path := filepath.Join("testdata", base+".html")
		file, err := os.Create(path)
		if err != nil {
			t.Errorf("Error creating file: %v", err)
		}
		defer file.Close()

		_, err = file.Write(bytes)
		if err != nil {
			t.Errorf("Error writing to file: %v", err)
		}
	}

	_, err = fetcher.Fetch(context.Background(), "not a valid url")
	if err == nil {
		t.Error("Failed to return error for invalid url")
	}
}

func TestFetchBookmark(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	fetcher, err := NewFetcher()
	assert.NilError(t, err)

	for idx, url := range urls {
		bookmark, err := fetcher.FetchBookmark(context.Background(), url)
		if err != nil {
			t.Errorf("Failed to fetch %s", url)
		}
		assert.Equal(t, bookmark.Title, titles[idx])
	}

	_, err = fetcher.FetchBookmark(context.Background(), "not a valid url")
	if err == nil {
		t.Error("Failed to return error for invalid url")
	}
}
