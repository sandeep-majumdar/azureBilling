package utils

import (
	"io"
	"os"
)

type FileSystem interface {
	Create(string) (io.WriteCloser, error)
	Open(string) (io.ReadCloser, error)
}

type LocalFS struct{}

func (LocalFS) Create(name string) (io.WriteCloser, error) { return os.Create(name) }
func (LocalFS) Open(name string) (io.ReadCloser, error)    { return os.Open(name) }
