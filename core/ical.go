package core

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/apognu/gocal"
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

func DownloadIcal(icalUrl string) error {
	return downloadFile("./temp/cal.ics", icalUrl)
}

func ParseIcal() (*gocal.Gocal, error) {
	f, err := os.Open("./temp/cal.ics")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour) // now to next 12 months

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	err = c.Parse()
	if err != nil {
		return nil, err
	}

	return c, nil
}
