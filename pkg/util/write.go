package util

import (
	"os"
)

func WriteToFile(path string, content *[]string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range *content {
		if _, err = f.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}
