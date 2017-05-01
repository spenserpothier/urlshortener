package main

type MyUrl struct {
	Title       string `json:"Title"`
	ExpandedUrl string `json:"url" schema:"url"`
	ShortUrl    string `json:"ShortUrl"`
}
