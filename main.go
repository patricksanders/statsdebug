package main

import (
	"fmt"
)

type Counter map[string]int

func (c Counter) increment(key string) error {
	c[key]++
	fmt.Printf("Stat %s has value %d\n", key, c[key])
	return nil
}

func (c Counter) get(key string) int {
	return c[key]
}

var cache *Counter

func main() {
	cache = &Counter{}
	go httpListener()
	statsdListener()
}
