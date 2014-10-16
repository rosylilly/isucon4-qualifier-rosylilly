package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	port = flag.Int("port", 8080, "the port to listen on")
)

const (
	defaultPath     = uint8('?')
	rootPath        = uint8('/')
	loginPath       = uint8('l')
	mypagePath      = uint8('m')
	reportPath      = uint8('r')
	imagesPath      = uint8('i')
	stylesheetsPath = uint8('s')
)

var (
	AYATAKA_MODE = false
)

type MainHandler struct {
}

func (m *MainHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start)
		log.Printf("%s \t%s\t%s", req.Method, duration, req.URL.Path)
	}()
	rawReq := req

	path := defaultPath
	if len(rawReq.URL.Path) > 1 {
		path = rawReq.URL.Path[1]
	} else {
		path = rootPath
	}

	switch path {
	case imagesPath:
		static(rw, rawReq)
	case stylesheetsPath:
		static(rw, rawReq)
	case loginPath:
		login(rw, rawReq)
	case mypagePath:
		myPage(rw, rawReq)
	case reportPath:
		report(rw, rawReq)
	case rootPath:
		topPage(rw, rawReq)
	default:
		http.NotFound(rw, rawReq)
	}
}

func main() {
	flag.Parse()

	portStr := strconv.Itoa(*port)

	server := &http.Server{
		Addr:    ":" + portStr,
		Handler: new(MainHandler),
	}
	log.Fatal(server.ListenAndServe())
}
