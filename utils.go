package main

import (
	"os"
)

func initTempFolder(tempPath string) error {
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err = os.Mkdir(tempPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
