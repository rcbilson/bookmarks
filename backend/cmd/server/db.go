package main

import (
	"context"
	"database/sql"
	"os"
	"unicode"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"
)

type Db interface {
	Close()
	Hit(ctx context.Context, url string) error
	SetFavorite(ctx context.Context, url string, isFavorite bool) error
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
	_, err := os.Stat(dbfile)
	if err != nil {
		_, err = os.Create(dbfile)
		if err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	schemaVersion := 0
	row := db.QueryRow("SELECT schemaVersion FROM metadata WHERE id = 0")
	_ = row.Scan(&schemaVersion)

	err = applySchema(db, schemaVersion)
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

	err = applySchema(db, 0)
	if err != nil {
		return nil, err
	}

	return &DbContext{db}, err
}

func applySchema(db *sql.DB, lastVersion int) error {
	for _, sql := range schema[lastVersion:] {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}
	_, err := db.Exec(`INSERT INTO metadata (id, schemaVersion) VALUES (0, @version)
						ON CONFLICT DO UPDATE SET schemaVersion = @version`,
		sql.Named("version", len(schema)))
	if err != nil {
		return err
	}
	return nil
}

func (ctx *DbContext) Close() {
	ctx.db.Close()
}

// Marks a bookmark as being frequently accessed
func (dbctx *DbContext) Hit(ctx context.Context, url string) error {
	_, err := dbctx.db.Exec("UPDATE bookmarks SET hitCount = hitCount + 1 WHERE url = ?", url)
	return err
}

// Marks a bookmark as being a favorite, or not
func (dbctx *DbContext) SetFavorite(ctx context.Context, url string, isFavorite bool) error {
	favorite := 0
	if isFavorite {
		favorite = 1
	}
	_, err := dbctx.db.ExecContext(ctx, "UPDATE bookmarks SET favorite = ? WHERE url = ?", favorite, url)
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

func scanBookmarkList(rows *sql.Rows) (bookmarkList, error) {
	var result bookmarkList

	for rows.Next() {
		var r bookmarkEntry
		var favorite int
		err := rows.Scan(&r.Title, &r.Url, &favorite)
		if err != nil {
			return nil, err
		}
		if favorite == 1 {
			r.IsFavorite = true
		} else {
			r.IsFavorite = false
		}
		result = append(result, r)
	}
	return result, nil
}

// Returns the most recently-accessed bookmarks
func (dbctx *DbContext) Recents(ctx context.Context, count int) (bookmarkList, error) {
	rows, err := dbctx.db.QueryContext(ctx, `SELECT title, url, favorite FROM bookmarks WHERE title != '""' ORDER BY lastAccess DESC LIMIT ?`, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBookmarkList(rows)
}

// Returns the most frequently-accessed bookmarks
func (dbctx *DbContext) Favorites(ctx context.Context, count int) (bookmarkList, error) {
	rows, err := dbctx.db.QueryContext(ctx, `SELECT title, url, favorite FROM bookmarks WHERE title != '""' AND favorite = 1 ORDER BY hitCount DESC LIMIT ?`, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBookmarkList(rows)
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
