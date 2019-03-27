package main

import (
	"sync"

	"github.com/inconshreveable/log15"
)

type Counter map[string]int

var log = log15.New("module", "statsdebug")
var lock = sync.RWMutex{}

func (c Counter) increment(key string) {
	lock.Lock()
	defer lock.Unlock()
	c[key]++
}

func (c Counter) get(key string) int {
	lock.RLock()
	defer lock.RUnlock()
	return c[key]
}

var counter *Counter

func init() {
}
func main() {
	log.Info("starting")
	counter = &Counter{}
	go statsdListener()
	serve()
}
