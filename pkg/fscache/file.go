package fscache

import (
	"bytes"
	"errors"
	"os"
)

type file struct {
	bytes.Reader
	stat    os.FileInfo
	statErr error
}

func (f *file) Close() error {
	return nil
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("Cached fs unable to readdir")
}

func (f *file) Stat() (os.FileInfo, error) {
	return f.stat, f.statErr
}
