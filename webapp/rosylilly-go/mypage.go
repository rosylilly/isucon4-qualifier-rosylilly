package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var (
	myPageTemplate *template.Template
)

func myPage(res http.ResponseWriter, req *http.Request) {
	var user *User = nil

	if cookie, err := req.Cookie("u"); err == nil {
		userId, _ := strconv.Atoi(cookie.Value)

		if userId > 0 && userId <= len(Users) {
			user = Users[userId-1]
		}
	}

	if user == nil {
		res.Header().Add("Set-Cookie", fmt.Sprintf("_=%d", LoginMust))
	}

	res.WriteHeader(200)
	myPageTemplate.Execute(res, user)
}
