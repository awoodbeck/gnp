package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Item struct {
	Title     string `xml:"title"`
	URL       string `xml:"link"`
	Published string `xml:"pubDate"`
}

type RSS struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
	entityTag string
}

func (r RSS) Items() []Item {
	items := make([]Item, len(r.Channel.Items))
	copy(items, r.Channel.Items)

	return items
}

func (r *RSS) ParseURL(ctx context.Context, u string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	if r.entityTag != "" {
		req.Header.Add("ETag", r.entityTag)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNotModified: // no-op
	case http.StatusOK:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()

		err = xml.Unmarshal(b, r)
		if err != nil {
			return err
		}

		r.entityTag = resp.Header.Get("ETag")
	default:
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return nil
}
