// Package irc provides an client-side implementation of the IRC protocol
// compatible with most modern IRC servers.
package irc

// TODO RFC-compliance

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
