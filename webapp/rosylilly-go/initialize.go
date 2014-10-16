package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type FlashMessage struct {
	Flash string
}

func Flash(message string) *FlashMessage {
	return &FlashMessage{message}
}

func ExecTemplate(tmpl *template.Template, data interface{}) []byte {
	buf := bytes.NewBufferString("")

	err := tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func init() {
	AYATAKA_MODE = os.Getenv("AYATAKA") == "1"

	fp, err := os.Create("rosylilly-go.pid")
	if err != nil {
		panic(err)
	}
	io.WriteString(fp, fmt.Sprintf("%d", os.Getpid()))
	fp.Close()

	var topPageTmpl *template.Template
	if AYATAKA_MODE {
		topPageTmpl = template.Must(template.ParseFiles("ayataka/topPage.html"))
		myPageTemplate = template.Must(template.ParseFiles("ayataka/myPage.html"))
	} else {
		topPageTmpl = template.Must(template.ParseFiles("html/topPage.html"))
		myPageTemplate = template.Must(template.ParseFiles("html/myPage.html"))
	}

	NoramlTopPage = ExecTemplate(topPageTmpl, nil)
	BannedTopPage = ExecTemplate(topPageTmpl, Flash("You're banned."))
	LockedTopPage = ExecTemplate(topPageTmpl, Flash("This account is locked."))
	WrongTopPage = ExecTemplate(topPageTmpl, Flash("Wrong username or password"))
	MustLoginTopPage = ExecTemplate(topPageTmpl, Flash("You must be logged in"))

	initalize()
	loadDump("dump_ips.json", &IPs)
	loadDump("dump_users.json", &Users)
	MakeUsersIndex()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for true {
			<-c
			log.Printf("Receive HUP")
			initalize()
		}
	}()
}

func loadDump(path string, data interface{}) {
	fp, err := os.Open(path)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}

	json.Unmarshal(bytes, data)
}

func initalize() {
	IPs = map[string]*IP{}
	loadUsersByTSV()
}

func loadUsersByTSV() {
	Users = []*User{}
	UserLoginIndex = map[string]*User{}

	usersTSV, err := os.Open("../../sql/dummy_users.tsv")
	if err != nil {
		panic(err)
	}
	userUsedTSV, err := os.Open("../../sql/dummy_users_used.tsv")
	if err != nil {
		panic(err)
	}

	failureCounts := map[string]uint32{}

	reader := csv.NewReader(userUsedTSV)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		failed, _ := strconv.Atoi(record[2])
		failureCounts[record[1]] = uint32(failed)
	}

	reader = csv.NewReader(usersTSV)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		uid, _ := strconv.Atoi(record[0])
		user := NewUser(uid, record[1], record[2])

		failed := uint32(0)
		if definedFailed, ok := failureCounts[user.Login]; ok {
			failed = definedFailed
		}
		user.FailCount = failed

		Users = append(Users, user)
		UserLoginIndex[user.Login] = user
	}
}
