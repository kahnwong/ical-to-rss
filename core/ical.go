package core

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/apognu/gocal"

	"github.com/rs/zerolog"
)

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadIcal(icalUrl string, logger zerolog.Logger) {
	err := downloadFile("./temp/cal.ics", icalUrl)
	if err != nil {
		logger.Error().Err(err).Msg("Error downloading ical file")
	}
}

func ParseIcal(logger zerolog.Logger) *gocal.Gocal {
	f, _ := os.Open("./temp/cal.ics")
	defer f.Close()

	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour) // now to next 12 months

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	err := c.Parse()
	if err != nil {
		logger.Error().Err(err).Msg("Error downloading ical file")
	}

	return c
}
