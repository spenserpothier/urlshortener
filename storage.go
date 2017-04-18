package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const DB_VERSION = 1

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
        Id INTEGER PRIMARY KEY NOT NULL,
        Title TEXT,
        ExpandedUrl TEXT,
        ShortUrl TEXT,
        NumberOfClicks INTEGER NOT NULL,
        InsertedDatetime DATETIME
);
PRAGMA user_version=1;
`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
	log.Printf("Table Successfully created")
}

func CheckForDBUpdates(db *sql.DB) {
	sql_check_version := `PRAGMA user_version`
	rows, err := db.Query(sql_check_version)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var v int
	for rows.Next() {
		rows.Scan(&v)
	}
	log.Printf("Existing DB Version: %v\n", v)
	log.Printf("Current DB_VERSION = %v\n", DB_VERSION)
	if v < DB_VERSION {
		log.Printf("Updating DB")
		UpdateDB(db, v)
	}
}

func UpdateDB(db *sql.DB, v int) {

	if v == 1 {
		// Do nothing, leaving here to remember how to add updates to DB down the road
		//log.Printf("Updating DB to Version 2")
		//updateVersion := "PRAGMA user_version=2;"
		//db.Exec(updateVersion)
		//v++
	}
}

func StoreUrl(db *sql.DB, url MyUrl) {
	sql_additem := `
INSERT OR REPLACE INTO MyUrls(
       Title,
       ExpandedUrl,
       ShortUrl,
       NumberOfClicks,
       InsertedDatetime
) values (?, ?, ?, 0, CURRENT_TIMESTAMP)
`
	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(&url.Title, &url.ExpandedUrl, &url.ShortUrl)
	log.Printf("Stored URL: %s with ShortURL: https://r.spenser.io/%s", url.ExpandedUrl, url.ShortUrl)
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
	sql_update_counter := `
UPDATE MyUrls SET NumberOfClicks = NumberOfClicks + 1
WHERE ShortUrl = ?
`
	stmt, _ = db.Prepare(sql_update_counter)
	_, err3 := stmt.Exec(s)
	if err3 != nil {
		panic(err)
	}
	return item
}

func GetAllUrls(db *sql.DB) []MyUrl {
	sql_find := `
SELECT Title, ShortUrl, ExpandedUrl FROM MyUrls
`
	stmt, _ := db.Prepare(sql_find)
	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []MyUrl

	for rows.Next() {
		item := MyUrl{}
		err2 := rows.Scan(&item.Title, &item.ShortUrl, &item.ExpandedUrl)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}
