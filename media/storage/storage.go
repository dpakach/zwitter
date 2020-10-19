package storage

import (
	"io"
	"os"
)

type Storage interface {
	Save(path string, file io.Reader) error
	Get(path string) (*os.File, error)
}
