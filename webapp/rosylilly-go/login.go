package main

import (
	"fmt"
	"net/http"
)

const (
	LoginSuccess = iota
	LoginBannedIP
	LoginLockedUser
	LoginWrongPassword
	LoginMust
)

func login(res http.ResponseWriter, req *http.Request) {
	login := req.PostFormValue("login")
	password := req.PostFormValue("password")

	remoteAddr := req.RemoteAddr
	if xForwardedFor := req.Header.Get("X-Forwarded-For"); len(xForwardedFor) > 0 {
		remoteAddr = xForwardedFor
	}

	result := LoginSuccess

	redirectPath := "/mypage"

	user, found := UserLoginIndex[login]
	ip := GetIP(remoteAddr)

	defer func() {
		if result != LoginSuccess {
			if user != nil {
				user.Fail()
			}
			if ip != nil {
				ip.Fail()
			}
		}
	}()

	if ip.IsBanned() {
		result = LoginBannedIP
	}

	if found {
		if result == LoginSuccess && user.IsLocked() {
			result = LoginLockedUser
		}

		if result == LoginSuccess && user.Password != password {
			result = LoginWrongPassword
		}
	} else {
		result = LoginWrongPassword
	}

	if result != LoginSuccess {
		redirectPath = "/"
	}

	if result != LoginSuccess {
		res.Header().Add("Set-Cookie", fmt.Sprintf("_=%d", result))
	} else {
		user.Success(remoteAddr)
		ip.Success()
		res.Header().Add("Set-Cookie", fmt.Sprintf("u=%d", user.Id))
	}

	http.Redirect(res, req, redirectPath, 302)
}
