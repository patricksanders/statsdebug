package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

type Stat struct {
	Name       string
	Type       string
	Value      string
	SampleRate string
	Tags       []string
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
	log.Info("handling stat packet", "stat", s.Name)
	summary.add(s)
}

func parseStat(stat string) (*Stat, error) {
	s := Stat{}

	parts := strings.SplitN(stat, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid stat string")
	}
	s.Name = parts[0]

	parts = strings.SplitN(parts[1], "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid stat string")
	}
	s.Value = parts[0]

	if strings.Contains(parts[1], "|") {
		// Another pipe means a sample rate.
		if strings.Contains(parts[1], "#") {
			// We also have tags!
			parts = strings.SplitN(parts[1], "|", 2)
			s.Type = parts[0]
			parts = strings.SplitN(parts[1], "#", 2)
			s.SampleRate = parts[0]
			s.Tags = strings.Split(parts[1], ",")
		} else {
			// No hash, so no tags. Pull out the type and sample rate.
			parts = strings.SplitN(parts[1], "|", 2)
			s.Type = parts[0]
			s.SampleRate = parts[1]
		}
	} else {
		// No sample rate here.
		if strings.Contains(parts[1], "#") {
			// We have tags!
			parts = strings.SplitN(parts[1], "#", 2)
			s.Type = parts[0]
			s.Tags = strings.Split(parts[1], ",")
		} else {
			// No hash, so no tags. Pull out the type.
			s.Type = parts[1]
		}
	}

	// Sort tags
	sort.Strings(s.Tags)

	return &s, nil
}
