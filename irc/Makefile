include ${GOROOT}/src/Make.inc

TARG = irc
GOFILES = irc.go message.go reply.go util.go clientproto.go client.go bot.go

include ${GOROOT}/src/Make.pkg

fmt:
	gofmt -w *.go

