package main

import (
	"context"
	"database/sql"
	"unicode"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"
)

type Db interface {
	Close()
	Hit(ctx context.Context, url string) error
	Get(ctx context.Context, url string) (BookmarkData, bool)
	Recents(ctx context.Context, count int) (bookmarkList, error)
	Favorites(ctx context.Context, count int) (bookmarkList, error)
	Insert(ctx context.Context, url string, bookmark BookmarkData) error
	Search(ctx context.Context, pattern string) (bookmarkList, error)
}

type DbContext struct {
	db *sql.DB
}

func NewDb(dbfile string) (Db, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	return &DbContext{db}, nil
}

func NewTestDb() (*DbContext, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
CREATE TABLE bookmarks (
  url text primary key,
  title text,
  lastAccess datetime,
  hitCount integer
);
CREATE VIRTUAL TABLE fts USING fts5(
  url UNINDEXED,
  title,
  content='bookmarks',
  prefix='1 2 3',
  tokenize='porter unicode61'
);
CREATE TRIGGER bookmarks_ai AFTER INSERT ON bookmarks BEGIN
  INSERT INTO fts(rowid, url, title) VALUES (new.rowid, new.url, new.title);
END;
        `)
	if err != nil {
		return nil, err
	}

	return &DbContext{db}, err
}

func (ctx *DbContext) Close() {
	ctx.db.Close()
}

// Marks a bookmark as being frequently accessed
func (dbctx *DbContext) Hit(ctx context.Context, url string) error {
	_, err := dbctx.db.Exec("UPDATE bookmarks SET hitCount = hitCount + 1 WHERE url = ?", url)
	return err
}

// Returns a bookmark title if one exists in the database
func (dbctx *DbContext) Get(ctx context.Context, url string) (BookmarkData, bool) {
	row := dbctx.db.QueryRowContext(ctx, "SELECT title FROM bookmarks WHERE url = ?", url)
	var title string
	err := row.Scan(&title)
	if err != nil {
		return BookmarkData{}, false
	}
	_, _ = dbctx.db.Exec("UPDATE bookmarks SET lastAccess = datetime('now') WHERE url = ?", url)
	return BookmarkData{Title: title}, true
}

// Returns the most recently-accessed bookmarks
func (dbctx *DbContext) Recents(ctx context.Context, count int) (bookmarkList, error) {
	rows, err := dbctx.db.QueryContext(ctx, `SELECT title, url FROM bookmarks WHERE title != '""' ORDER BY lastAccess DESC LIMIT ?`, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result bookmarkList

	for rows.Next() {
		var r bookmarkEntry
		err := rows.Scan(&r.Title, &r.Url)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

// Returns the most frequently-accessed bookmarks
func (dbctx *DbContext) Favorites(ctx context.Context, count int) (bookmarkList, error) {
	rows, err := dbctx.db.QueryContext(ctx, `SELECT title, url FROM bookmarks WHERE title != '""' ORDER BY hitCount DESC LIMIT ?`, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result bookmarkList

	for rows.Next() {
		var r bookmarkEntry
		err := rows.Scan(&r.Title, &r.Url)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

// Insert the bookmark title corresponding to the url into the database
func (dbctx *DbContext) Insert(ctx context.Context, url string, bookmark BookmarkData) error {
	_, err := dbctx.db.ExecContext(ctx, "INSERT INTO bookmarks (url, title, lastAccess, hitCount) VALUES (?, ?, datetime('now'), 0)", url, bookmark.Title)
	return err
}

// Search for bookmarks matching a pattern
func (dbctx *DbContext) Search(ctx context.Context, pattern string) (bookmarkList, error) {
	if pattern == "" {
		return nil, nil
	}
	// If the final token in the pattern is a letter, add a star to treat it as
	// a prefix query
	lastRune, _ := utf8.DecodeLastRuneInString(pattern)
	if unicode.IsLetter(lastRune) {
		pattern += "*"
	}
	rows, err := dbctx.db.QueryContext(ctx, "SELECT title, url FROM fts where fts MATCH ? ORDER BY rank", pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result bookmarkList

	for rows.Next() {
		var r bookmarkEntry
		err := rows.Scan(&r.Title, &r.Url)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
