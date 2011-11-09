package main

import (
	"flag"
	"log"
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
			handlePriv(conn, m.From, m.To, m.Text)
		default:
			log.Print("WARNING: unhandled message kind")
		}
	}
}

func handlePriv(c Client, from, to, text string) {
	if to[0] != '#' {
		log.Printf("Private message from %s: %s", from, text)
		return
	}
	c.PrivMsg(to, "Ditto")
}
