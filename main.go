package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/feeds"

	"github.com/apognu/gocal"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// download ical file
	icalUrl := os.Getenv("ICAL_URL")
	err := downloadFile("cal.ics", icalUrl)
	if err != nil {
		fmt.Println("Error downloading ical file")
		os.Exit(1)
	}

	// parse ical
	f, _ := os.Open("cal.ics")
	defer f.Close()

	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour) // now to next 12 months

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	err = c.Parse()
	if err != nil {
		fmt.Println("Cannot parse ical")
		os.Exit(1)
	}

	//for _, e := range c.Events {
	//	fmt.Printf("%s\n%s\n%s\n\n%s", e.Summary, e.Start, e.Organizer.Cn, e.Description)
	//	fmt.Println("---")
	//}

	// generate rss feed
	now := time.Now()
	feed := &feeds.Feed{
		Title:       os.Getenv("FEED_TITLE"),
		Link:        &feeds.Link{Href: os.Getenv("ICAL_URL")},
		Description: os.Getenv("FEED_DESCRIPTION"),
		Author:      &feeds.Author{Name: os.Getenv("FEED_AUTHOR_NAME"), Email: os.Getenv("FEED_AUTHOR_EMAIL")},
		Created:     now,
	}

	var feedItems []*feeds.Item
	for _, e := range c.Events {
		feedItems = append(feedItems, &feeds.Item{
			Title:       e.Summary,
			Link:        &feeds.Link{Href: ""},
			Description: e.Description,
			Author:      &feeds.Author{Name: e.Organizer.Cn, Email: e.Organizer.Cn},
			Created:     *e.Created,
		})
	}

	feed.Items = feedItems
	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rss)
}
