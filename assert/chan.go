package assert

import (
	"errors"
	"fmt"
)

// Ensure that a value is immediately sent on the specified channel.
func ChanSend(c chan<- interface{}, v interface{}) error {
	select {
	case c <- v:
		return nil
	default:
		return errors.New("sending on channel failed")
	}
}

// Ensure that a value is received on the specified channel.
func ChanRecv(c <-chan interface{}) (interface{}, error) {
	select {
	case v := <-c:
		return v, nil
	default:
		return nil, errors.New("receiving on channel failed")
	}
}

// Ensure that the specified value is received on the channel.
func ChanRecvVal(c <-chan interface{}, v interface{}) error {
	if recv, err := ChanRecv(c); err == nil {
		if recv != v {
			return fmt.Errorf("%v != %v", recv, v)
		} else {
			return nil
		}
	} else {
		return err
	}
}

// Ensure that the channel is closed.
func ChanClosed(c <-chan interface{}) error {
	select {
	case _, ok := <-c:
		if !ok {
			return nil
		}
	default:
	}
	return errors.New("channel is not closed")
}
