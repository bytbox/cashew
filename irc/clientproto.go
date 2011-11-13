package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

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

// Low-level method to connect to server - normal clients should not need this.
// Use Connect() instead.
func Dial(server string) (conn *Client, err error) {
	nconn, err := net.Dial("tcp", server)
	if err != nil {
		return
	}

	conn = new(Client)
	conn.Conn = nconn
	conn.serverName = server
	conn.server = make(chan ServerMessage, serverMsgBufSize)

	// spawn the connection reader
	go func() {
		r, err := bufio.NewReaderSize(conn, maxlinesize)
		if err != nil {
			panic(err)
		}

		line, beFalse, err := r.ReadLine()
		for err == nil && !beFalse {
			conn.server <- parseServerMessage(string(line))
			line, beFalse, err = r.ReadLine()
		}
		if beFalse {
			log.Fatal("Line too long")
		} else {
			log.Fatal(err) // TODO handle me better
		}
	}()
	return
}

// Low-level method to send PASS command - normal clients should not need this.
func (c *Client) Pass(pass string) {
	fmt.Fprintf(c, "PASS %s\n", pass)
}

// Low-level method to send USER command - normal clients should not need this.
func (c *Client) User(user, host, server, name string) {
	fmt.Fprintf(c, "USER %s %s %s :%s\n", user, host, server, name)
}

// Low-level method to send NICK command - normal clients should not need this.
func (c *Client) Nick(nick string) {
	fmt.Fprintf(c, "NICK %s\n", nick)
}

// Low-level method to join the specified channel. This does not modify the
// Client's internal channel tracking, and so should not be used by most
// clients.
func (c *Client) Join(ch string) {
	// TODO don't join until 001 is received
	log.Print("JOIN ", ch)
	fmt.Fprintf(c, "JOIN %s\n", ch)
}

// Join the specified channels. Equivalent to calling Join for each channel in the given slice.
func (c *Client) JoinChannels(chs []string) {
	for _, ch := range chs {
		c.Join(ch)
	}
}

// Low-level method to leave the specified channel. This does not modify the
// Client's internal channel tracking, and so should not be used by most
// clients.
func (c *Client) Part(ch string) {
	log.Print("PART ", ch)
	fmt.Fprintf(c, "PART %s\n", ch)
}

// Low-level method to send QUIT to the server.
func (c *Client) Quit(msg string) {
	log.Print("QUIT :", msg)
	fmt.Fprintf(c, "QUIT :%s\n", msg)
}

// Low-level method to send a private message (PRIVMSG).
func (c *Client) PrivMsg(to, msg string) {
	fmt.Fprintf(c, "PRIVMSG %s :%s\n", to, msg)
}
