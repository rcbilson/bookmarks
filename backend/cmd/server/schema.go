package main

var schema = []string{
	// version 1
	`
CREATE TABLE metadata (
  id integer primary key,
  schemaVersion integer
);

create table if not exists bookmarks (
  url text primary key,
  title text,
  lastAccess datetime,
  hitCount integer
);

DROP TABLE IF EXISTS fts;

CREATE VIRTUAL TABLE fts USING fts5(
  url UNINDEXED,
  title,
  content='bookmarks',
  prefix='1 2 3',
  tokenize='porter unicode61'
);

-- Triggers to keep the FTS index up to date.
DROP TRIGGER IF EXISTS bookmarks_ai;
CREATE TRIGGER bookmarks_ai AFTER INSERT ON bookmarks BEGIN
  INSERT INTO fts(rowid, url, title) VALUES (new.rowid, new.url, new.title);
END;

DROP TRIGGER IF EXISTS bookmarks_ad;
CREATE TRIGGER bookmarks_ad AFTER DELETE ON bookmarks BEGIN
  INSERT INTO fts(fts, rowid, url, title) VALUES('delete', old.rowid, old.url, old.title);
END;

DROP TRIGGER IF EXISTS bookmarks_au;
CREATE TRIGGER bookmarks_au AFTER UPDATE ON bookmarks BEGIN
  INSERT INTO fts(fts, rowid, url, title) VALUES('delete', old.rowid, old.url, old.title);
  INSERT INTO fts(rowid, url, title) VALUES (new.rowid, new.url, new.title);
END;

INSERT INTO fts(fts) VALUES('rebuild');
	`,

	// version 2
	`
ALTER TABLE bookmarks ADD COLUMN favorite integer DEFAULT 0;
	`,
}
