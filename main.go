package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
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
	CheckForDBUpdates(db)

	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/add", AddHandler)
	r.HandleFunc("/list", ListHandler)
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

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newUrl MyUrl
		r.ParseForm()
		log.Printf("form: %v", r.PostForm) //
		decoder := schema.NewDecoder()
		err := decoder.Decode(&newUrl, r.PostForm)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error %v: ", err)
			io.WriteString(w, "Error")
			return
		}
		if len(newUrl.ShortUrl) == 0 {
			newUrl.ShortUrl = randSeq(7)
		}
		StoreUrl(db, newUrl)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "New url: <a href=\"https://r.spenser.io/%s\">https://r.spenser.io/%s</a>", newUrl.ShortUrl, newUrl.ShortUrl)
	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/add.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		p := &struct {
			Title string
			Url   string
		}{
			r.URL.Query().Get("title"),
			r.URL.Query().Get("url"),
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, p)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Method not supported")
		return
	}

}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/list.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, GetAllUrls(db))
}

func ShortenedHandler(w http.ResponseWriter, r *http.Request) {
	pa := mux.Vars(r)
	p := pa["key"]
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
