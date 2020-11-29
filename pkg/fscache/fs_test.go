package fscache

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dir http.FileSystem = http.Dir(".")

func TestMemFS_open_not_existed_file(t *testing.T) {
	fs, terminate := NewFSCache(http.Dir("."))
	defer terminate()

	_, err := fs.Open("./no_file")
	_, rerr := dir.Open("./no_file")

	assert.Equal(t, rerr, err)
	assert.NotNil(t, err)
}

func TestMemFS_release_file_with_ttl(t *testing.T) {
	ioutil.WriteFile("release_file_with_ttl", []byte("ok"), 0777)
	defer os.Remove("release_file_with_ttl")

	fs, terminate := NewFSCache(http.Dir("."))
	defer terminate()
	fs.SetTtl(1)

	fmt.Println("o1")
	file, err := fs.Open("./release_file_with_ttl")

	fmt.Println("o2")
	assertData(t, "ok", file)
	fmt.Println("o3")
	assert.Nil(t, err)

	fmt.Println("o4")
	os.Remove("release_file_with_ttl")
	time.Sleep(time.Second * 6)

	fmt.Println("o4")
	file, err = fs.Open("./release_file_with_ttl")

	assert.NotNil(t, err)
}

func TestMemFS_open_existed_file(t *testing.T) {
	ioutil.WriteFile("open_existed_file", []byte("ok"), 0777)
	defer os.Remove("open_existed_file")

	fmt.Println("sdfasd")
	fs, terminate := NewFSCache(http.Dir("."))
	defer terminate()

	file, err := fs.Open("./open_existed_file")
	rfile, rerr := dir.Open("./open_existed_file")

	assert.Equal(t, rerr, err)
	assert.Nil(t, err)

	data := make([]byte, 2)
	rfile.Read(data)
	assert.Equal(t, "ok", string(data))

	assertData(t, "ok", file)
}

func TestMemFS_open_removed_file(t *testing.T) {
	ioutil.WriteFile("open_removed_file", []byte("ok"), 0777)

	fmt.Println("sdfasd")
	fs, terminate := NewFSCache(http.Dir("."))
	defer terminate()

	file, err := fs.Open("./open_removed_file")
	assert.Nil(t, err)
	os.Remove("open_removed_file")

	file, err = fs.Open("./open_removed_file")
	assert.Nil(t, err)

	assertData(t, "ok", file)
}

func BenchmarkMemFS_open_existed_file(b *testing.B) {
	ioutil.WriteFile("bench_open_existed_file", []byte("ok"), 0777)
	defer os.Remove("bench_open_existed_file")

	fs, terminate := NewFSCache(http.Dir("."))
	defer terminate()

	for i := 0; i < b.N; i++ {

		file, err := fs.Open("./bench_open_existed_file")
		assert.Nil(b, err)

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		assert.Equal(b, "ok", buf.String())
	}
}

func assertData(t *testing.T, expected string, file io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	assert.Equal(t, expected, buf.String())
}
