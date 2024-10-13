package main

import (
	"fmt"
	"os"
	"time"

	"github.com/apognu/gocal"
	"github.com/gorilla/feeds"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		//logger.Fatal.Err(err).Msg("") // [TODO]
	} else {
		// write to file
		err = os.WriteFile("feed.rss", []byte(rss), 0644)
	}

	fmt.Println(rss) // [TODO] remove

	// api server
	// entrypoint
	mode := os.Getenv("MODE")
	listenAddress := ""
	isPrettyLog := false
	switch mode {
	case "PRODUCTION":
		listenAddress = ":3000"
	case "DEVELOPMENT":
		listenAddress = "localhost:3000"
		isPrettyLog = true
	default:
		log.Fatal().Msg("Listen address is not set")
	}

	// app
	app := fiber.New()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if isPrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	// 60 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        60,
	}))

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to subsonic-widgets api")
	})

	app.Static("/feed.rss", "./feed.rss")

	if err := app.Listen(listenAddress); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}
}
