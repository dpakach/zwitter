package storage

import (
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	BasePath string
}

func NewLocal(BasePath string) (*Local, error) {
	p, err := filepath.Abs(BasePath)
	if err != nil {
		return nil, err
	}
	return &Local{p}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	fullPath := filepath.Join(l.BasePath, path)

	dir := filepath.Dir(fullPath)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, contents)
	if err != nil {
		return err
	}
	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	fullPath := filepath.Join(l.BasePath, path)
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f, nil
}
