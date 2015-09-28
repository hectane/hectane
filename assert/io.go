package assert

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

// Ensure that reading from the specified reader results in the specified data.
func Read(r io.Reader, data []byte) error {
	if b, err := ioutil.ReadAll(r); err == nil {
		if reflect.DeepEqual(b, data) {
			return nil
		} else {
			return errors.New(fmt.Sprintf("%v != %v", b, data))
		}
	} else {
		return err
	}
}
