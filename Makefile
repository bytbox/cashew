include ${GOROOT}/src/Make.inc

TARG = gabon
GOFILES = gabon.go irc/irc.go irc/reply.go irc/message.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w *.go

