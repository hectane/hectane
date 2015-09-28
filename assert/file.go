package assert

import (
	"errors"
	"os"
)

// Ensure file is in the specified state or return an error.
func FileState(filename string, exists bool) error {
	if _, err := os.Stat(filename); err == nil {
		if exists {
			return nil
		} else {
			return errors.New("file exists")
		}
	} else {
		if os.IsNotExist(err) {
			if exists {
				return errors.New("file does not exist")
			} else {
				return nil
			}
		} else {
			return err
		}
	}
}
