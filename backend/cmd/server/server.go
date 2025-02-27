package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Port         int    `default:"9000"`
	FrontendPath string `default:"/home/richard/src/bookmark/frontend/dist"`
	DbFile       string `default:"/home/richard/src/bookmark/data/bookmark.db"`
}

var spec specification

func main() {
	err := envconfig.Process("bookmarkserver", &spec)
	if err != nil {
		log.Fatal("error reading environment variables:", err)
	}

	db, err := NewDb(spec.DbFile)
	if err != nil {
		log.Fatal("error initializing database interface:", err)
	}
	defer db.Close()

	fetcher, err := NewFetcher()
	if err != nil {
		log.Fatal("error initializing fetcher:", err)
	}

	handler(db, fetcher, spec.Port, spec.FrontendPath)
}
