package fscache

import (
	"bytes"
	"net/http"
	"os"
)

type cachedFile struct {
	ttl  int
	data []byte
	stat os.FileInfo

	readErr error
	openErr error
	statErr error
}

func (f *cachedFile) file() (http.File, error) {
	return &file{
		Reader:  *bytes.NewReader(f.data),
		stat:    f.stat,
		statErr: f.statErr,
	}, f.openErr
}
