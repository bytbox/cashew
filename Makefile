include ${GOROOT}/src/Make.inc

TARG = gabon
GOFILES = gabon.go irc.go reply.go message.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w *.go
