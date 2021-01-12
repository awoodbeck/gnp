package main

import (
	"context"

	"github.com/awoodbeck/gnp/ch14/feed"
	"github.com/aws/aws-lambda-go/lambda"
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

func main() {
	lambda.Start(LatestXKCD)
}

func LatestXKCD(ctx context.Context, req EventRequest) (
	EventResponse, error) {
	resp := EventResponse{Title: "xkcd.com", URL: "https://xkcd.com/"}

	if err := rssFeed.ParseURL(ctx, feedURL); err != nil {
		return resp, err
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

	return resp, nil
}
