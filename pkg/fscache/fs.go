package fscache

import (
	"bytes"
	"errors"
	"net/http"
	"sync"
	"time"
)

func NewFSCache(fs http.FileSystem) (*FSCache, terminate) {
	cache := &FSCache{
		ttl:   60,
		files: make(map[string]cachedFile),
		mu:    sync.Mutex{},
		fs:    fs,
	}

	quit := make(chan bool)
	ticker := time.NewTicker(time.Second * 5)

	go func() {
		defer close(quit)
		defer ticker.Stop()

		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				cache.gc()
			}
		}
	}()

	return cache, func() { quit <- true }
}

type FSCache struct {
	ttl   int
	mu    sync.Mutex
	files map[string]cachedFile
	fs    http.FileSystem
}

func (cache *FSCache) Open(name string) (http.File, error) {
	file, ok := cache.files[name]
	if ok {
		file.ttl++
		if file.ttl > cache.ttl {
			file.ttl = cache.ttl
		}
		return file.file()
	}

	file = cachedFile{ttl: cache.ttl}

	rf, err := cache.fs.Open(name)
	if err != nil {
		file.openErr = err
		cache.mu.Lock()
		cache.files[name] = file
		cache.mu.Unlock()

		return file.file()
	}
	defer rf.Close()

	stat, statErr := rf.Stat()
	file.stat = stat
	file.statErr = statErr

	buf := new(bytes.Buffer)
	_, readErr := buf.ReadFrom(rf)
	file.data = buf.Bytes()
	file.readErr = readErr

	cache.mu.Lock()
	cache.files[name] = file
	cache.mu.Unlock()

	return file.file()
}

func (cache *FSCache) SetTtl(ttl int) error {
	if ttl < 0 {
		return errors.New("ttl should be greater then zero")
	}
	cache.ttl = ttl

	return nil
}

func (cache *FSCache) gc() {
	for name, file := range cache.files {
		file.ttl -= 5
		if file.ttl < 1 {
			cache.mu.Lock()
			delete(cache.files, name)
			cache.mu.Unlock()
		}
	}
}

type terminate = func()
