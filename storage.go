package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//	"log"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTable(db *sql.DB) {
	sql_table := `
CREATE TABLE IF NOT EXISTS MyUrls(
        Title TEXT,
        ExpandedUrl TEXT,
        ShortUrl TEXT,
        InsertedDatetime DATETIME
);
`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func StoreUrl(db *sql.DB, url MyUrl) {
	sql_additem := `
INSERT OR REPLACE INTO MyUrls(
       Title,
       ExpandedUrl,
       ShortUrl,
       InsertedDatetime
) values (?, ?, ?, CURRENT_TIMESTAMP)
`
	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(&url.Title, &url.ExpandedUrl, &url.ShortUrl)
	if err2 != nil {
		panic(err2)
	}
}

func FindUrl(db *sql.DB, s string) MyUrl {
	sql_find := `
SELECT Title, ShortUrl, ExpandedUrl FROM MyUrls
WHERE ShortUrl = ? LIMIT 1
`
	stmt, _ := db.Prepare(sql_find)
	rows, err := stmt.Query(s)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	item := MyUrl{}

	for rows.Next() {
		err2 := rows.Scan(&item.Title, &item.ShortUrl, &item.ExpandedUrl)
		if err2 != nil {
			panic(err2)
		}
	}
	return item
}
