package main

import (
	"encoding/json"
	"net/http"
)

func report(res http.ResponseWriter, req *http.Request) {
	bannedIPs := []string{}
	lockedUsers := []string{}

	for addr, ip := range IPs {
		if ip.IsBanned() {
			bannedIPs = append(bannedIPs, addr)
		}
	}

	for _, user := range Users {
		if user.IsLocked() {
			lockedUsers = append(lockedUsers, user.Login)
		}
	}

	bytes, err := json.Marshal(map[string][]string{
		"banned_ips":   bannedIPs,
		"locked_users": lockedUsers,
	})
	if err != nil {
		bytes = []byte{}
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(bytes)
}
