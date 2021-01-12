package gcp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/awoodbeck/gnp/ch14/feed"
)

var (
	rssFeed feed.RSS
	feedURL = "https://xkcd.com/rss.xml"
)

type EventRequest struct {
	Previous bool `json:"previous"`
}

type EventResponse struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Published string `json:"published"`
}

func LatestXKCD(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	resp := EventResponse{Title: "xkcd.com", URL: "https://xkcd.com/"}

	defer func() {
		w.Header().Set("Content-Type", "application/json")
		out, _ := json.Marshal(&resp)
		_, _ = w.Write(out)
	}()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("decoding request: %v", err)
		return
	}

	if err := rssFeed.ParseURL(r.Context(), feedURL); err != nil {
		log.Printf("parsing feed: %v:", err)
		return
	}

	switch items := rssFeed.Items(); {
	case req.Previous && len(items) > 1:
		resp.Title = items[1].Title
		resp.URL = items[1].URL
		resp.Published = items[1].Published
	case len(items) > 0:
		resp.Title = items[0].Title
		resp.URL = items[0].URL
		resp.Published = items[0].Published
	}
}
