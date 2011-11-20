include ${GOROOT}/src/Make.inc

TARG = cashew
GOFILES = cashew.go util.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w *.go

