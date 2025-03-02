package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gotest.tools/assert"
)

type mockFetcher struct {
}

func (*mockFetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	return []byte("<html><head><title>title for " + url + "</title></head></html>"), nil
}

func (*mockFetcher) FetchBookmark(_ context.Context, url string) (BookmarkData, error) {
	return BookmarkData{Title: "title for " + url + "</title></head></html>"}, nil
}

type titleStruct struct {
	Title string `json:"title"`
}

type bookmarkListEntryStruct struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type bookmarkListStruct []bookmarkListEntryStruct

var testFetcher = &mockFetcher{}

func addTest(t *testing.T, db Db, reqUrl string) {
	v := url.Values{}
	v.Add("url", reqUrl)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/add?url=%s", v.Encode()), nil)
	w := httptest.NewRecorder()
	add(db, testFetcher)(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func listTest(t *testing.T, handler func(http.ResponseWriter, *http.Request), reqName string, reqCount int, expCount int, resultList *bookmarkListStruct) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s?count=%d", reqName, reqCount), nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	var bookmarkList bookmarkListStruct
	err := json.NewDecoder(resp.Body).Decode(&bookmarkList)
	assert.NilError(t, err)
	assert.Equal(t, expCount, len(bookmarkList))
	if resultList != nil {
		*resultList = bookmarkList
	}
}

func searchTest(t *testing.T, db Db, pattern string, expCount int) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/search?q=%s", url.QueryEscape(pattern)), nil)
	w := httptest.NewRecorder()
	search(db)(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	var bookmarkList bookmarkListStruct
	err := json.NewDecoder(resp.Body).Decode(&bookmarkList)
	assert.NilError(t, err)
	assert.Equal(t, expCount, len(bookmarkList))
}

func hitTest(t *testing.T, db Db, urlstr string) {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/hit?url=%s", url.QueryEscape(urlstr)), nil)
	w := httptest.NewRecorder()
	hit(db)(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

// TODO: test something other than the happy path
func TestHandlers(t *testing.T) {
	db, err := NewTestDb()
	assert.NilError(t, err)

	// basic add request
	addTest(t, db, urls[0])

	// repeating test should produce same result but hit db
	addTest(t, db, urls[0])

	// set up a second title in the db
	addTest(t, db, urls[1])

	// ask for five recents, expect two
	listTest(t, fetchRecents(db), "recent", 5, 2, nil)

	// ask for one recent, expect one
	listTest(t, fetchRecents(db), "recent", 1, 1, nil)

	// ask for one favorite, expect one
	listTest(t, fetchFavorites(db), "favorite", 1, 1, nil)

	// ask for five favorites, expect two
	var resultList bookmarkListStruct
	listTest(t, fetchFavorites(db), "favorite", 5, 2, &resultList)

	// hit whichever was reported second
	hitTest(t, db, resultList[1].Url)

	// ask for the favorites after the hit, second should now be first
	var newResultList bookmarkListStruct
	listTest(t, fetchFavorites(db), "favorite", 2, 2, &newResultList)
	assert.Equal(t, resultList[1].Title, newResultList[0].Title)

	// should have one search hit
	searchTest(t, db, "seriouseats", 1)

	// prefix should be allowed
	searchTest(t, db, "serious", 1)

	// should have two search hits
	searchTest(t, db, "www", 2)

	// should have no search hits
	searchTest(t, db, "foo", 0)
}
