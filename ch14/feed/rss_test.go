package feed

import (
	"context"
	"testing"
)

func TestRSSParseURL(t *testing.T) {
	rss := new(RSS)
	err := rss.ParseURL(context.Background(), "https://xkcd.com/rss.xml")
	if err != nil {
		t.Fatal(err)
	}

	if rss.entityTag == "" {
		t.Error("unexpected empty ETag")
	}

	if len(rss.Items()) == 0 {
		t.Error("empty Items after successful ParseURL call")
	}

	etag := rss.entityTag
	err = rss.ParseURL(context.Background(), "https://xkcd.com/rss.xml")
	if err != nil {
		t.Fatal(err)
	}

	if etag != rss.entityTag {
		t.Log("new ETag on second request")
	}

	if !t.Failed() {
		t.Logf("%#v", rss)
	}
}
