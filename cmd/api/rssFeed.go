package main

import (
	"net/http"
)

// this is my transport layer so all the structs can go here
// no need for any service layer or storage layer
// we are jsut getting some data = transport only
// this is not like a mailing service which we can use in our transport layer
// you can hit the rss with postman and click the preview to understand the structure
// also we are leaving some un-neccessary fields here

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"` 
	Description string `xml:"description"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

func (app *application) GetFeedHandler(w http.ResponseWriter, r *http.Request) {



}
