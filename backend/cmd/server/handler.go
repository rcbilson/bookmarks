package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type bookmarkEntry struct {
	Title      string `json:"title"`
	Url        string `json:"url"`
	IsFavorite bool   `json:"isFavorite"`
}

type bookmarkList []bookmarkEntry

func handler(db Db, fetcher Fetcher, port int, frontendPath string) {
	// Handle the api routes in the backend
	http.Handle("POST /api/add", http.HandlerFunc(add(db, fetcher)))
	http.Handle("GET /api/recents", http.HandlerFunc(fetchRecents(db)))
	http.Handle("GET /api/favorites", http.HandlerFunc(fetchFavorites(db)))
	http.Handle("GET /api/search", http.HandlerFunc(search(db)))
	http.Handle("POST /api/hit", http.HandlerFunc(hit(db)))
	http.Handle("POST /api/setFavorite", http.HandlerFunc(setFavorite(db)))
	// bundled assets and static resources
	http.Handle("GET /assets/", http.FileServer(http.Dir(frontendPath)))
	http.Handle("GET /static/", http.FileServer(http.Dir(frontendPath)))
	// For other requests, serve up the frontend code
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", frontendPath))
	})
	log.Println("server listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func logError(w http.ResponseWriter, msg string, code int) {
	log.Printf("%d %s", code, msg)
	http.Error(w, msg, code)
}

func search(db Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok {
			logError(w, "No search terms provided", http.StatusBadRequest)
			return
		}
		list, err := db.Search(r.Context(), query[0])
		if err != nil {
			logError(w, fmt.Sprintf("Error fetching recent bookmarks: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)
		w.Header().Set("Content-Type", "application/json")
	}
}

func fetchRecents(db Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		count := 5
		countStr, ok := r.URL.Query()["count"]
		if ok {
			count, err = strconv.Atoi(countStr[0])
			if err != nil {
				logError(w, fmt.Sprintf("Invalid count specification: %s", countStr[0]), http.StatusBadRequest)
				return
			}
		}
		recentList, err := db.Recents(r.Context(), count)
		if err != nil {
			logError(w, fmt.Sprintf("Error fetching recent bookmarks: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(recentList)
		w.Header().Set("Content-Type", "application/json")
	}
}

func fetchFavorites(db Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		count := 5
		countStr, ok := r.URL.Query()["count"]
		if ok {
			count, err = strconv.Atoi(countStr[0])
			if err != nil {
				logError(w, fmt.Sprintf("Invalid count specification: %s", countStr[0]), http.StatusBadRequest)
				return
			}
		}
		recentList, err := db.Favorites(r.Context(), count)
		if err != nil {
			logError(w, fmt.Sprintf("Error fetching favorite bookmarks: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(recentList)
		w.Header().Set("Content-Type", "application/json")
	}
}

func hit(db Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := r.URL.Query()["url"]
		if !ok {
			logError(w, "No url provided", http.StatusBadRequest)
			return
		}
		err := db.Hit(r.Context(), url[0])
		if err != nil {
			logError(w, fmt.Sprintf("Error updating database: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func setFavorite(db Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := r.URL.Query()["url"]
		if !ok {
			logError(w, "No url provided", http.StatusBadRequest)
			return
		}
		isFavorite, ok := r.URL.Query()["isFavorite"]
		if !ok {
			logError(w, "No url provided", http.StatusBadRequest)
			return
		}
		favorite := false
		if isFavorite[0] == "true" {
			favorite = true
		} else if isFavorite[0] != "false" {
			logError(w, "Expected true/false for isFavorite", http.StatusBadRequest)
			return
		}
		err := db.SetFavorite(r.Context(), url[0], favorite)
		if err != nil {
			logError(w, fmt.Sprintf("Error updating database: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func add(db Db, fetcher Fetcher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()

		urls, ok := r.URL.Query()["url"]
		if !ok {
			logError(w, fmt.Sprintf("No url provided in request %v", r.URL), http.StatusBadRequest)
			return
		}
		url := urls[0]
		doUpdate := false
		bookmarkData, ok := db.Get(ctx, url)
		if !ok {
			log.Println("fetching bookmark", url)
			doUpdate = true
			bookmarkData, err = fetcher.FetchBookmark(ctx, url)
			if err != nil {
				logError(w, fmt.Sprintf("Error retrieving site: %v", err), http.StatusBadRequest)
				return
			}
		}
		if doUpdate {
			err = db.Insert(ctx, url, bookmarkData)
			if err != nil {
				log.Printf("Error inserting into db: %v", err)
			}
		}
	}
}
