package main

import (
	"flag"
	"log"
)

var server = "irc.freenode.net:6667"
var channels = []string{
	//"#go-nuts",
}

func main() {
	flag.Parse()

	log.Printf("Connecting to %s", server)
	conn, err := Connect(server, "gabon", "The Go Nut")
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	for _, c := range channels {
		conn.Join(c)
	}

	for {
		var m ServerMessage
		select {
		// TODO allow input on stdin
		case m = <-conn.server:
			handleMessage(m)
		}
	}
}

// TODO this should be shifted into irc.go
func handleMessage(m ServerMessage) {
	switch m.Code {
	case "NOTICE":
		log.Printf("NOTICE %s", m.Raw)
	case "PRIVMSG":
	case RPL_MOTD:
	case RPL_MOTDSTART:
	case RPL_ENDOFMOTD:
	case RPL_WELCOME:
	case RPL_YOURHOST:
	case RPL_CREATED:
	case RPL_MYINFO:
	case RPL_BOUNCE:
	case RPL_LUSERCLIENT:
	case RPL_LUSEROP:
	case RPL_LUSERUNKNOWN:
	case RPL_LUSERCHANNELS:
	case RPL_LUSERME:
	case "MODE":
	default:
		log.Printf("Unhandled message %s\t\t%s", m.Code, m.Raw)
	}
}
