package main

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

var (
	BasePath = "./public"
)

var (
	staticFiles  map[string][]byte
	contentTypes map[string]string
)

func init() {
	staticFiles = map[string][]byte{}
	contentTypes = map[string]string{}

	files := []string{
		"/images/isucon-bank.png",
		"/stylesheets/bootflat.min.css",
		"/stylesheets/bootstrap.min.css",
		"/stylesheets/isucon-bank.css",
	}

	for _, path := range files {
		fp, err := os.Open(BasePath + path)
		if err != nil {
			panic(err)
		}
		bytes, err := ioutil.ReadAll(fp)
		if err != nil {
			panic(err)
		}
		staticFiles[path] = bytes
		contentTypes[path] = mime.TypeByExtension(filepath.Ext(path))
	}
}

func static(res http.ResponseWriter, req *http.Request) {
	path := filepath.Clean(filepath.FromSlash(req.URL.Path))

	res.Header().Add("Content-Type", contentTypes[path])
	res.WriteHeader(200)
	res.Write(staticFiles[path])
}
