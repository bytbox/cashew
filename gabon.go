package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	. "github.com/bytbox/gabon/irc"
)

var server = "irc.freenode.net:6667"
var channels = []string{
	//"#go-nuts",
	"#go-bots",
}

func main() {
	flag.Parse()

	log.Printf("Connecting to %s", server)
	conn, err := Connect(server, "gabon", "The Go Nut")
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	conn.JoinChannels(channels)

	messages := conn.Listen()

	for m := range messages {
		switch m.Kind {
		case MSG_NOTICE:
			log.Printf("NOTICE\t%s", m.Text)
		case MSG_PRIVMSG:
			handlePriv(conn, getNick(m.From), m.To, m.Text)
		default:
			log.Print("WARNING: unhandled message kind")
		}
	}
}

func getNick(s string) string {
	a := strings.SplitN(s, "!", 2)
	return a[0]
}

func handlePriv(c *Client, from, to, text string) {
	replyTo := to
	if to[0] != '#' { // not sent to a channel - reply directly to user
		replyTo = from
	} else {
		// exit if this message isn't for us. TODO match ^! ?
		if r, _ := regexp.MatchString("gabon[,:] ", text); !r {
			return
		} else {
			_, text, _ = nextField(text)
		}
	}
	r, e := getReply(text)
	if r != "" {
		if e {
			r = fmt.Sprintf("%s, %s", from, r)
		}
		c.PrivMsg(replyTo, r)
	}
}

func getReply(text string) (string, bool) {
	f, rest, _ := nextField(text)
	if len(f) < 2 {
		return "", true
	}
	switch f[0] {
	case '@':
		r, e := getReply(rest)
		if !e {
			return f[1:] + ", " + r, false
		}
		return r, true
	}
	return "I don't understand "+f, true
}
