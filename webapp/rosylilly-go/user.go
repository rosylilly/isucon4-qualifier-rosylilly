package main

import (
	"sync"
	"time"
)

type User struct {
	*sync.Mutex

	Id              int       `json:"i"`
	Login           string    `json:"l"`
	Password        string    `json:"p"`
	LastLoginedAt   time.Time `json:"lt"`
	LastLoginedIP   string    `json:"li"`
	LatestLoginedAt time.Time `json:"llt"`
	LatestLoginedIP string    `json:"lli"`
	FailCount       uint32    `json:"f"`
}

var (
	Users          []*User
	UserLoginIndex map[string]*User
)

func NewUser(id int, login, password string) *User {
	return &User{
		Mutex:    new(sync.Mutex),
		Id:       id,
		Login:    login,
		Password: password,
	}
}

func MakeUsersIndex() {
	for _, user := range Users {
		UserLoginIndex[user.Login] = user
	}
}

func (u *User) IsLocked() bool {
	u.Lock()
	defer u.Unlock()

	return u.FailCount >= 3
}

func (u *User) Fail() {
	u.Lock()
	defer u.Unlock()
	u.FailCount++
}

func (u *User) Success(ip string) {
	if u.IsLocked() {
		return
	}

	u.Lock()
	u.LastLoginedAt = u.LatestLoginedAt
	u.LastLoginedIP = u.LatestLoginedIP
	u.LatestLoginedAt = time.Now()
	u.LatestLoginedIP = ip
	if len(u.LastLoginedIP) == 0 {
		u.LastLoginedAt = u.LatestLoginedAt
		u.LastLoginedIP = u.LatestLoginedIP
	}
	u.FailCount = 0
	u.Unlock()
}
