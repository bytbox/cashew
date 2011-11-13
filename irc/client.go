package irc

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	net.Conn
	serverName string
	server     chan ServerMessage
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
			fmt.Fprintf(c, "PONG :%s\n", sm.To)
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
