package main

import (
	"sync"
	"sync/atomic"
)

type IP struct {
	FailCount uint32
}

var (
	ipLock *sync.Mutex
	IPs    map[string]*IP
)

func init() {
	ipLock = new(sync.Mutex)
}

func GetIP(addr string) *IP {
	ipLock.Lock()
	if ip, ok := IPs[addr]; ok {
		ipLock.Unlock()
		return ip
	}

	ip := &IP{0}
	IPs[addr] = ip
	ipLock.Unlock()

	return ip
}

func (i *IP) IsBanned() bool {
	return atomic.LoadUint32(&i.FailCount) >= 10
}

func (i *IP) Fail() {
	atomic.AddUint32(&i.FailCount, 1)
}

func (i *IP) Success() {
	if i.IsBanned() {
		return
	}

	atomic.StoreUint32(&i.FailCount, 0)
}
