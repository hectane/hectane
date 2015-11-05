package cfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)
	var (
		c = Config{}
		f = filepath.Join(d, "config.json")
	)
	if err := c.Save(f); err != nil {
		t.Fatal(err)
	}
}
