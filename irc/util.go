package irc

import (
	"log"
	"os"
	"os/user"
	"strings"
)

func nextField(line string) (string, string, bool) {
	fs := strings.SplitN(line, " ", 2)
	if len(fs) > 1 {
		return fs[0], fs[1], true
	}
	return fs[0], "", false
}

func getUname() string {
	uid := os.Getuid()
	uname := "go"
	u, e := user.LookupId(uid)
	if e != nil {
		log.Print("WARNING: user.Lookupid: ", e)
	} else {
		uname = u.Username
	}
	return uname
}

func getHostname() string {
	n, e := os.Hostname()
	if e != nil {
		log.Print("WARNING: os.Hostname: ", e)
		n = "*"
	}
	return n
}

func getServername() string {
	return "*"
}
