package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	NoramlTopPage    []byte
	BannedTopPage    []byte
	LockedTopPage    []byte
	WrongTopPage     []byte
	MustLoginTopPage []byte
)

var (
	ExpireTime         = time.Unix(0, 0)
	ExpireTimeString   = ExpireTime.Format(time.UnixDate)
	RemoveCookieString = fmt.Sprintf("_=; expires=%s;", ExpireTimeString)
)

func topPage(res http.ResponseWriter, req *http.Request) {
	html := NoramlTopPage

	rmCookie := false
	if cookie, err := req.Cookie("_"); err == nil {
		switch cookie.Value {
		case "1":
			html = BannedTopPage
		case "2":
			html = LockedTopPage
		case "3":
			html = WrongTopPage
		default:
			html = MustLoginTopPage
		}
		rmCookie = true
	}

	if rmCookie {
		res.Header().Add("Set-Cookie", RemoveCookieString)
	}

	res.WriteHeader(200)
	res.Write(html)
}
