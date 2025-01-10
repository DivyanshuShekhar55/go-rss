package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
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
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// structures for json response ...
type FeedResponse struct {
	Title    string    `json:"title"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	PubDate     string `json:"pubDate"`
}

func (app *application) GetFeedHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("https://dev.to/rss")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Couldn't send request")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		writeJSONError(w, resp.StatusCode, "Failed to fetch rss feed")
		return
	}

	// read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to read response body")
		return
	}

	// parse xml
	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "faailed to parse rss feed")
		return
	}

	// convert to json format
	response := FeedResponse{
		Title:    rss.Channel.Title,
		Articles: make([]Article, 0),
		// items will be appended later
	}

	// convert rss items to article struct :
	for _, item := range rss.Channel.Items {
		article := Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PubDate:     item.PubDate,
		}

		response.Articles = append(response.Articles, article)
	}

	// send json response ...
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

}
