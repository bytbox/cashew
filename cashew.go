package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	. "github.com/bytbox/cashew/irc"
)

var server = "irc.freenode.net:6667"
var channels = []string{
	//"#go-nuts",
	"#go-bots",
}

func main() {
	flag.Parse()

	log.Printf("Connecting to %s", server)
	conn, err := Connect(server, "cashew", "The Go Nut")
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	conn.JoinChannels(channels)

	messages := conn.Listen()

	// catch signals
	go func() {
		for s := range signal.Incoming {
			log.Print("Got signal: ", s.String())
			// TODO part from all channels
			conn.Quit("Going nuts!")
			os.Exit(1)
		}
	}()

	// TODO recover from a panic
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
		if r, _ := regexp.MatchString("cashew[,:] ", text); !r {
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

type Command func(string) (string, bool)

var commands = map[string]Command{
	"learn": func(text string) (string, bool) {
		_, _, g := nextField(text)
		if !g {
			return "I need more than that!", true
		}
		return "I would have learned if I knew how...", true
	},
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
	case '!':
		f = f[1:]
	}
	// look for f in command index
	c, ok := commands[f]
	if ok {
		return c(rest)
	}
	return "I don't understand "+f, true
}
