package main

import (
	"fmt"
	"net"
	"strings"
)

type Stat struct {
	Name  string
	Type  string
	Value string
	Tags  []string
}

func statsdListener() {
	addr, _ := net.ResolveUDPAddr("udp", ":8125")
	sock, _ := net.ListenUDP("udp", addr)

	for {
		buf := make([]byte, 1024)
		rlen, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		go handleStat(string(buf[0:rlen]))
	}
}

func handleStat(stat string) {
	s, _ := parseStat(stat)
	cache.increment(s.Name)
}

func parseStat(stat string) (*Stat, error) {
	s := Stat{}
	parts := strings.SplitN(stat, ":", 2)
	s.Name = parts[0]
	parts = strings.SplitN(parts[1], "|", 2)
	s.Value = parts[0]
	parts = strings.SplitN(parts[1], "#", 2)
	s.Type = parts[0]
	s.Tags = strings.Split(parts[1], ",")

	return &s, nil
}
