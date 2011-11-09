package main

const (
	MSG_NOTICE = iota
	MSG_PRIVMSG
)

type Message struct {
	Kind int
	From string
	To   string
	Text string
}

func (m Message) IsPriv() bool {
	return m.Kind == MSG_PRIVMSG
}
