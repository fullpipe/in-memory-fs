package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/fullpipe/memfs/pkg/fscache"
)

var etagSeed int = rand.Int()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appRoot := os.Getenv("APP_ROOT")
	if appRoot == "" {
		appRoot = "/"
	}

	fs, terminate := fscache.NewFSCache(http.Dir("./app"))
	defer terminate()

	handler := http.FileServer(fs)
	handler = gziphandler.GzipHandler(handler)
	handler = indexFile(handler)
	handler = httpCache(handler)
	handler = onlyGetRequests(handler)
	handler = http.StripPrefix(appRoot, handler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.Handle("/", handler)

	log.Fatal(srv.ListenAndServe().Error())
}

func onlyGetRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func indexFile(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if filepath.Ext(r.URL.Path) == "" {
			r.URL.Path = "/"
		}

		h.ServeHTTP(w, r)
	})
}

func httpCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path
		e := `"` + key + "-" + strconv.Itoa(etagSeed) + `"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age=2592000") // 30 days

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
