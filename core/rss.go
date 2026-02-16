package core

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/apognu/gocal"
	"github.com/gorilla/feeds"
)

func GenerateRss(c *gocal.Gocal) (string, error) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       os.Getenv("FEED_TITLE"),
		Link:        &feeds.Link{Href: os.Getenv("ICAL_URL")},
		Description: os.Getenv("FEED_DESCRIPTION"),
		Author:      &feeds.Author{Name: os.Getenv("FEED_AUTHOR_NAME"), Email: os.Getenv("FEED_AUTHOR_EMAIL")},
		Created:     now,
	}

	var feedItems []*feeds.Item
	slices.Reverse(c.Events)
	for _, e := range c.Events {
		feedItems = append(feedItems, &feeds.Item{
			Title:       e.Summary,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://th.techcal.dev/%s", e.RawStart.Value)},
			Description: strings.ReplaceAll(e.Description, "\\n", "<br>"),
			Author:      &feeds.Author{Name: "", Email: ""},
			Created:     *e.Created,
		})
	}

	feed.Items = feedItems
	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}
