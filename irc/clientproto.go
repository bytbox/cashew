package irc

import (
	"fmt"
	"log"
)

// Low-level method to send PASS command - normal clients should not need this.
func (c *Client) Pass(pass string) {
	fmt.Fprintf(c.out, "PASS %s\n", pass)
}

// Low-level method to send USER command - normal clients should not need this.
func (c *Client) User(user, host, server, name string) {
	fmt.Fprintf(c.out, "USER %s %s %s :%s\n", user, host, server, name)
}

// Low-level method to send NICK command - normal clients should not need this.
func (c *Client) Nick(nick string) {
	fmt.Fprintf(c.out, "NICK %s\n", nick)
	c.nick = nick // TODO fix possible race condition
}


// Low-level method to join the specified channel. This does not modify the
// Client's internal channel tracking, and so should not be used by most
// clients.
func (c *Client) Join(ch string) {
	// TODO don't join until 001 is received
	log.Print("JOIN ", ch)
	fmt.Fprintf(c.out, "JOIN %s\n", ch)
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
	fmt.Fprintf(c.out, "PART %s\n", ch)
}

// Low-level method to send QUIT to the server.
func (c *Client) Quit(msg string) {
	log.Print("QUIT :", msg)
	fmt.Fprintf(c.out, "QUIT :%s\n", msg)
}

// Low-level method to send a private message (PRIVMSG).
func (c *Client) PrivMsg(to, msg string) {
	fmt.Fprintf(c.out, "PRIVMSG %s :%s\n", to, msg)
}
