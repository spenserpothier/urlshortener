package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var db *sql.DB

func main() {
	const dbpath = "myurls.db"

	db = InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	r := mux.NewRouter()
	r.HandleFunc("/add", AddHandler)
	r.HandleFunc("/list", HomeHandler)
	r.HandleFunc("/{key}", ShortenedHandler)
	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("certificates/development.cer", "certificates/development.key"))
}

var myShortenedUrls = make(map[string]*MyUrl)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Method not supported")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Error")
		return
	}
	var newUrl MyUrl
	json.Unmarshal(body, &newUrl)

	defer r.Body.Close()
	if len(newUrl.ShortUrl) == 0 {
		newUrl.ShortUrl = randSeq(7)
	}
	StoreUrl(db, newUrl)
	log.Printf("Stored URL: %s\n", newUrl.ExpandedUrl)
	//myShortenedUrls[s] = &newUrl
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "New url: <a href=\"https://r.spenser.io/%s\">https://r.spenser.io/%s</a>", newUrl.ShortUrl, newUrl.ShortUrl)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var msg string
	for k, v := range myShortenedUrls {
		msg += fmt.Sprintf("<a href=\"https://r.spenser.io/%s\">https://r.spenser.io/%s</a> - %s - %s<br>", k, k, v.Title, v.ExpandedUrl)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, msg)
}

func ShortenedHandler(w http.ResponseWriter, r *http.Request) {
	//p := r.URL.EscapedPath()
	pa := mux.Vars(r)
	p := pa["key"]
	log.Printf("Escaped path: %s", p)
	s := FindUrl(db, p)
	if s.ExpandedUrl != "" {
		http.Redirect(w, r, s.ExpandedUrl, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
