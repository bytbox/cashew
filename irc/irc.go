package irc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"strings"
)

// Size of internal buffer for reading. 4kB should be enough to cover even the
// poorest implementation's idea of a longest permissible line.
const maxlinesize = 4096

// Maximum length of a line we output - anything else will be wrapped
const lineLength = 80

// Size of the ServerMessage buffer
const serverMsgBufSize = 10

// A message sent by the server
type ServerMessage struct {
	From string
	Code string
	To   string

	// The raw content of the message, excluding From, Code, and To
	Raw string

	// The full faw content of the message
	Full string

	// All other fields.
	Fields []string
}

type Client struct {
	connection net.Conn
	serverName string
	nick       string
	server     chan ServerMessage
	out        io.Writer
}

func nextField(line string) (string, string, bool) {
	fs := strings.SplitN(line, " ", 2)
	if len(fs) > 1 {
		return fs[0], fs[1], true
	}
	return fs[0], "", false
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

// Low-level method to connect to server - normal clients should not need this.
// Use Connect() instead.
func Dial(server string) (conn *Client, err error) {
	nconn, err := net.Dial("tcp", server)
	if err != nil {
		return
	}

	conn = new(Client)
	conn.connection = nconn
	conn.serverName = server
	conn.server = make(chan ServerMessage, serverMsgBufSize)
	conn.out = conn.connection

	// spawn the connection reader
	go func() {
		r, err := bufio.NewReaderSize(conn.connection, maxlinesize)
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

func (c *Client) Pass(pass string) {
	fmt.Fprintf(c.out, "PASS %s\n", pass)
}

// Low-level method to send USER command - normal clients should not need this
func (c *Client) User(user, host, server, name string) {
	fmt.Fprintf(c.out, "USER %s %s %s :%s\n", user, host, server, name)
}

// Low-level method to send NICK command - normal clients should not need this.
func (c *Client) Nick(nick string) {
	fmt.Fprintf(c.out, "NICK %s\n", nick)
	c.nick = nick // TODO fix possible race condition
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
			// TODO allow input on stdin
			case m = <-c.server:
				handleMessage(m)
			}
		}
	}()
	return ch
}

// Join the specified channel
func (c *Client) Join(ch string) {
	// TODO don't join until 001 is received
	log.Print("JOIN ", ch)
	fmt.Fprintf(c.out, "JOIN %s\n", ch)
}

// Join the specified channels
func (c *Client) JoinChannels(chs []string) {
	for _, ch := range chs {
		c.Join(ch)
	}
}

// Leave the specified channel
func (c *Client) Part(ch string) {
	log.Print("PART ", ch)
	fmt.Fprintf(c.out, "PART %s\n", ch)
}

func (c *Client) Quit(msg string) {
	log.Print("QUIT :", msg)
	fmt.Fprintf(c.out, "QUIT :%s\n", msg)
}

func (c *Client) PrivMsg(to, msg string) {
	fmt.Fprintf(c.out, "PRIVMSG %s :%s\n", to, msg)
}
