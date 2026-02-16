package main

import (
	"net/http"
	"os"

	"github.com/kahnwong/ical-to-rss/core"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func FeedHandler(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get ical
		if err := initTempFolder("./temp"); err != nil {
			logger.Error().Err(err).Msg("Error creating temp directory")
			c.String(http.StatusInternalServerError, "Error creating temp directory")
			return
		}

		icalUrl := os.Getenv("ICAL_URL")
		if err := core.DownloadIcal(icalUrl); err != nil {
			logger.Error().Err(err).Msg("Error downloading ical file")
			c.String(http.StatusInternalServerError, "Error downloading ical file")
			return
		}

		calendar, err := core.ParseIcal()
		if err != nil {
			logger.Error().Err(err).Msg("Error parsing ical file")
			c.String(http.StatusInternalServerError, "Error parsing ical file")
			return
		}

		// generate rss
		rss, err := core.GenerateRss(calendar)
		if err != nil {
			logger.Error().Err(err).Msg("Error generating RSS")
			c.String(http.StatusInternalServerError, "Error generating RSS")
			return
		}

		// serve response
		c.Data(http.StatusOK, "application/xml", []byte(rss))
	}
}
