package irc

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	connection net.Conn
	serverName string
	nick       string
	server     chan ServerMessage
	out        io.Writer
}

func parseServerMessage(line string) (m ServerMessage) {
	m.Full = line
	m.From, line, _ = nextField(line)
	if m.From[0] != ':' {
		m.Code = m.From
		m.From = ""
	} else {
		m.From = m.From[1:]
		m.Code, line, _ = nextField(line)
	}
	m.To, line, _ = nextField(line)
	m.Raw = line
	m.Fields = make([]string, 0)

	switch m.Code {
	case "NOTICE":
	case "PRIVMSG":
		m.Fields = append(m.Fields, line[1:])
	case RPL_BOUNCE:
	default:
		// fill in variable fields
		var f string
		f, line, b := nextField(line)
		for b {
			if f[0] == ':' {
				if f[1] == '-' {
					break
				}
				// read until f[len(f)-1] is a semicolon
				for b && len(f) > 0 && f[len(f)-1] != ':' {
					f, line, b = nextField(line)
				}
			} else {
				m.Fields = append(m.Fields, f)
			}
			f, line, b = nextField(line)
		}
	}
	return
}

func Connect(serverName, nick, realName string) (*Client, error) {
	conn, err := Dial(serverName)
	if err != nil {
		return conn, err
	}

	conn.Pass("notmeaningfull")
	conn.User(getUname(), getHostname(), getServername(), realName)
	conn.Nick(nick)
	return conn, nil
}

// Listen for messages coming in and return them on the returned channel. Also
// handles low-level information from the server correctly, making information
// available in the Client object as appropriate.
func (c *Client) Listen() <-chan Message {
	ch := make(chan Message)
	handleMessage := func(sm ServerMessage) {
		switch sm.Code {
		case "NOTICE":
			ch <- Message{
				Kind: MSG_NOTICE,
				From: sm.From,
				To:   sm.To,
				Text: sm.Raw[1:],
			}
		case "PRIVMSG":
			ch <- Message{
				Kind: MSG_PRIVMSG,
				From: sm.From,
				To:   sm.To,
				Text: sm.Raw[1:],
			}
		case "PING":
			fmt.Fprintf(c.out, "PONG :%s\n", sm.To)
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
			log.Printf("Unhandled message %s\t\t%s", sm.Code, sm.Full)
		}
	}
	go func() {
		for {
			var m ServerMessage
			select {
			case m = <-c.server:
				handleMessage(m)
			}
		}
	}()
	return ch
}
