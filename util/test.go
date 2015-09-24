package util

import (
	"errors"
	"fmt"
	"os"
)

// Ensure that a value is immediately sent on the specified channel.
func AssertChanSend(c chan<- interface{}, v interface{}) error {
	select {
	case c <- v:
		return nil
	default:
		return errors.New("sending on channel failed")
	}
}

// Ensure that a value is received on the specified channel.
func AssertChanRecv(c <-chan interface{}) (interface{}, error) {
	select {
	case v := <-c:
		return v, nil
	default:
		return nil, errors.New("receiving on channel failed")
	}
}

// Ensure that the specified value is received on the channel.
func AssertChanRecvVal(c <-chan interface{}, v interface{}) error {
	if recv, err := AssertChanRecv(c); err == nil {
		if recv != v {
			return errors.New(fmt.Sprintf("%v != %v", recv, v))
		} else {
			return nil
		}
	} else {
		return err
	}
}

// Ensure that the channel is closed.
func AssertChanClosed(c <-chan interface{}) error {
	select {
	case _, ok := <-c:
		if !ok {
			return nil
		}
	default:
	}
	return errors.New("channel is not closed")
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
