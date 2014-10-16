package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-c
		finalize()
		os.Exit(0)
	}()
}

func dumpJSON(data interface{}, path string) {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fp, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	fp.Write(bytes)
	fp.Close()
}

func finalize() {
	dumpJSON(Users, "dump_users.json")
	dumpJSON(IPs, "dump_ips.json")
}
