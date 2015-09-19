package util

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"sync"
)

var InvalidItem = errors.New("invalid item")

// Access broker for a key=>value store on disk with reference counting. All
// methods are safe to call from multiple goroutines.
type Filestore struct {
	sync.Mutex
	directory string
	items     map[string]int
}

// Determine the name of the index file on disk.
func (f *Filestore) indexFilename() string {
	return path.Join(f.directory, "index.json")
}

// Update the item map on disk.
func (f *Filestore) update() error {
	if w, err := os.OpenFile(f.indexFilename(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err == nil {
		if err := json.NewEncoder(w).Encode(f.items); err == nil {
			return w.Close()
		} else {
			return err
		}
	} else {
		return err
	}
}

// Create a new filestore. An attempt will be made to load the index from disk.
func NewFilestore(directory string) (*Filestore, error) {
	f := &Filestore{
		directory: directory,
		items:     make(map[string]int),
	}
	if r, err := os.Open(f.indexFilename()); err == nil {
		if err := json.NewDecoder(r).Decode(f.items); err != nil {
			return nil, err
		}
		if err := r.Close(); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return f, nil
}

// Obtain a writer for creating a new item.
func (f *Filestore) New(id string) (io.WriteCloser, error) {
	f.Lock()
	defer f.Unlock()
	if w, err := os.OpenFile(path.Join(f.directory, id), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err == nil {
		return w, nil
	} else {
		return nil, err
	}
}

// Obtain a reader for retrieving an exising item.
func (f *Filestore) Get(id string) (io.ReadCloser, error) {
	f.Lock()
	defer f.Unlock()
	if _, ok := f.items[id]; ok {
		return os.Open(path.Join(f.directory, id))
	} else {
		return nil, InvalidItem
	}
}

// Increase the reference count of the specified item.
func (f *Filestore) Add(id string) error {
	f.Lock()
	defer f.Unlock()
	if _, ok := f.items[id]; ok {
		f.items[id]++
	} else {
		f.items[id] = 1
	}
	return f.update()
}

// Decrease the reference count of the specified item. It will be deleted from
// disk when the reference count reaches zero.
func (f *Filestore) Release(id string) error {
	f.Lock()
	defer f.Unlock()
	if n, ok := f.items[id]; ok {
		if n == 1 {
			delete(f.items, id)
			if err := os.Remove(path.Join(f.directory, id)); err != nil {
				return err
			}
		} else {
			f.items[id]--
		}
		return f.update()
	} else {
		return InvalidItem
	}
}
