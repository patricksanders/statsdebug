package main

import (
	"github.com/inconshreveable/log15"
)

var (
	log = log15.New("module", "statsdebug")
	summary *Summary
)

func init() {
	summary = NewSummary()
}

func main() {
	log.Info("starting")
	go statsdListener()
	serve()
}
