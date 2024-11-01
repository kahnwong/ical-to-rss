package core

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/apognu/gocal"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog"
)

func GenerateRss(c *gocal.Gocal, logger zerolog.Logger) string {
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
			Description: e.Description,
			Author:      &feeds.Author{Name: "", Email: ""},
			Created:     *e.Created,
		})
	}

	feed.Items = feedItems
	rss, err := feed.ToRss()
	if err != nil {
		logger.Error().Err(err).Msg("Error generating RSS")
	}

	return rss
}
