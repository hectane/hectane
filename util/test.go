package util

import (
	"errors"
	"os"
)

// Ensure that a value is immediately sent on the specified channel.
func AssertChanSend(c chan interface{}, v interface{}) error {
	select {
	case c <- v:
		return nil
	default:
		return errors.New("sending on channel failed")
	}
}

// Ensure file is in the specified state or return an error.
func AssertFileState(filename string, exists bool) error {
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
