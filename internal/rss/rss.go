package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}
	client := &http.Client{}
	req.Header.Add("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching rss from %s: %w", feedURL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	decoder := xml.NewDecoder(resp.Body)
	var rss RSSFeed
	if err := decoder.Decode(&rss); err != nil {
		return nil, fmt.Errorf("error parsing rss feed: %w", err)
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for i := range len(rss.Channel.Item) {
		item := rss.Channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rss.Channel.Item[i] = item
	}

	return &rss, nil
}
