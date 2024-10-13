package main

import (
	"os"

	"github.com/rs/zerolog"
)

func initTempFolder(tempPath string, logger zerolog.Logger) {
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err = os.Mkdir(tempPath, 0755)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Error creating %s directory", tempPath)
		}
	}
}
