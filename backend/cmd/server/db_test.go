package main

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func setupTest(t *testing.T) *DbContext {
	db, err := NewTestDb()
	assert.NilError(t, err)
	t.Cleanup(db.Close)
	return db
}

func TestInsertGet(t *testing.T) {
	db := setupTest(t)
	ctx := context.Background()

	assert.NilError(t, db.Insert(ctx, "http://example.com", BookmarkData{Title: `bookmark`}))
	assert.Assert(t, nil != db.Insert(ctx, "http://example.com", BookmarkData{Title: `bookmark`}))
	bookmark, ok := db.Get(ctx, "http://example.com")
	assert.Assert(t, ok)
	assert.Equal(t, bookmark.Title, `bookmark`)
	_, ok = db.Get(ctx, "http://foo.com")
	assert.Assert(t, !ok)
}

func TestRecents(t *testing.T) {
	db := setupTest(t)
	ctx := context.Background()

	// set up two bookmarks
	assert.NilError(t, db.Insert(ctx, "http://example.com", BookmarkData{Title: `bookmark`}))
	assert.NilError(t, db.Insert(ctx, "http://example2.com", BookmarkData{Title: `bookmark2`}))
	assert.NilError(t, db.Insert(ctx, "http://example3.com", BookmarkData{Title: `""`}))

	// ask for 5, expect 2
	recents, err := db.Recents(ctx, 5)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(recents))
}

func TestFavorites(t *testing.T) {
	db := setupTest(t)
	ctx := context.Background()

	// set up two bookmarks
	assert.NilError(t, db.Insert(ctx, "http://example.com", BookmarkData{Title: `bookmark`}))
	assert.NilError(t, db.Insert(ctx, "http://example2.com", BookmarkData{Title: `bookmark2`}))

	// no favorites yet
	faves, err := db.Favorites(ctx, 5)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(faves))

	// set as favorites
	assert.NilError(t, db.SetFavorite(ctx, "http://example.com", true))
	assert.NilError(t, db.SetFavorite(ctx, "http://example2.com", true))

	// ask for 5, expect 2
	faves, err = db.Favorites(ctx, 5)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(faves))

	// hit the one in second place, it should come first
	secondPlace := faves[1].Url
	err = db.Hit(ctx, secondPlace)
	assert.NilError(t, err)
	newFaves, err := db.Favorites(ctx, 1)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(newFaves))
	assert.Equal(t, secondPlace, newFaves[0].Url)
}

func TestSearch(t *testing.T) {
	db := setupTest(t)
	ctx := context.Background()

	// set up two bookmarks
	assert.NilError(t, db.Insert(ctx, "http://example.com", BookmarkData{Title: `one two"}`}))
	assert.NilError(t, db.Insert(ctx, "http://example2.com", BookmarkData{Title: `one three"}`}))

	// expect 2
	results, err := db.Search(ctx, "one")
	assert.NilError(t, err)
	assert.Equal(t, 2, len(results))

	// expect 1
	results, err = db.Search(ctx, "one two")
	assert.NilError(t, err)
	assert.Equal(t, 1, len(results))

	// expect 0
	results, err = db.Search(ctx, "one two three")
	assert.NilError(t, err)
	assert.Equal(t, 0, len(results))

	// expect 1, auto prefix final token
	results, err = db.Search(ctx, "one thr")
	assert.NilError(t, err)
	assert.Equal(t, 1, len(results))

	// expect 1, phrase match
	results, err = db.Search(ctx, `"one three"`)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(results))

	// expect 0, no auto prefix
	results, err = db.Search(ctx, `"one thr"`)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(results))
}

func TestGetUpdatesLastAccessed(t *testing.T) {
	db := setupTest(t)
	ctx := context.Background()

	_, err := db.db.Exec(`INSERT INTO bookmarks (url, title, lastAccess) VALUES ('http://example.com', 'bookmark', '2016-03-29')`)
	assert.NilError(t, err)
	_, err = db.db.Exec(`INSERT INTO bookmarks (url, title, lastAccess) VALUES ('http://example2.com', 'bookmark2', '2016-03-30')`)
	assert.NilError(t, err)

	// example2 should be the first result
	recents, err := db.Recents(ctx, 1)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(recents))
	assert.Equal(t, "http://example2.com", recents[0].Url)
	assert.Equal(t, "bookmark2", recents[0].Title)

	// a Get on example should make it the first result
	_, ok := db.Get(ctx, "http://example.com")
	assert.Equal(t, true, ok)
	recents, err = db.Recents(ctx, 1)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(recents))
	assert.Equal(t, "http://example.com", recents[0].Url)
	assert.Equal(t, "bookmark", recents[0].Title)
}

func TestInsertUpdatesLastAccessed(t *testing.T) {
	db := setupTest(t)
	ctx := context.Background()

	_, err := db.db.Exec(`INSERT INTO bookmarks (url, title, lastAccess) VALUES ('http://example2.com', 'bookmark2', '2016-03-30')`)
	assert.NilError(t, err)

	// example2 should be the first result
	recents, err := db.Recents(ctx, 1)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(recents))
	assert.Equal(t, "http://example2.com", recents[0].Url)
	assert.Equal(t, "bookmark2", recents[0].Title)

	// a inserting example should make it the first result
	assert.NilError(t, db.Insert(ctx, "http://example.com", BookmarkData{Title: "bookmark"}))
	recents, err = db.Recents(ctx, 1)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(recents))
	assert.Equal(t, "http://example.com", recents[0].Url)
	assert.Equal(t, "bookmark", recents[0].Title)
}
