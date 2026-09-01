package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/kahnwong/ical-to-rss/core"
)

func FeedHandler(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		// get ical
		if err := initTempFolder("./temp"); err != nil {
			logger.Error("Error creating temp directory", "error", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating temp directory")
		}

		icalUrl := os.Getenv("ICAL_URL")
		if err := core.DownloadIcal(icalUrl); err != nil {
			logger.Error("Error downloading ical file", "error", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error downloading ical file")
		}

		calendar, err := core.ParseIcal()
		if err != nil {
			logger.Error("Error parsing ical file", "error", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error parsing ical file")
		}

		// generate rss
		rss, err := core.GenerateRss(calendar)
		if err != nil {
			logger.Error("Error generating RSS", "error", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error generating RSS")
		}

		// serve response
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)
		return c.SendString(rss)
	}
}
