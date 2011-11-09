include ${GOROOT}/src/Make.inc

TARG = gabon
GOFILES = gabon.go util.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w *.go

