package main

import (
	"os"
	"time"

	"github.com/kahnwong/ical-to-rss/core"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// init
	mode := os.Getenv("MODE")

	// ------------ api server ------------ //
	// entrypoint
	listenAddress := ""
	switch mode {
	case "PRODUCTION":
		listenAddress = ":3000"
	case "DEVELOPMENT":
		listenAddress = "localhost:3000"
	default:
		log.Fatal().Msg("Listen address is not set")
	}

	// app
	app := fiber.New()

	// set logger
	isPrettyLog := false
	if mode == "DEVELOPMENT" {
		isPrettyLog = true
	}
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
		return c.SendString("ICAL to RSS")
	})
	app.Get("/feed", func(c *fiber.Ctx) error {
		// get ical
		if err := initTempFolder("./temp"); err != nil {
			logger.Error().Err(err).Msg("Error creating temp directory")
			return err
		}

		icalUrl := os.Getenv("ICAL_URL")
		if err := core.DownloadIcal(icalUrl); err != nil {
			logger.Error().Err(err).Msg("Error downloading ical file")
			return err
		}

		calendar, err := core.ParseIcal()
		if err != nil {
			logger.Error().Err(err).Msg("Error parsing ical file")
			return err
		}

		// generate rss
		rss, err := core.GenerateRss(calendar)
		if err != nil {
			logger.Error().Err(err).Msg("Error generating RSS")
			return err
		}

		// serve rersponse
		c.Type("xml")
		_, err = c.Write([]byte(rss))

		return err
	})

	// start server
	if err := app.Listen(listenAddress); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}
}
