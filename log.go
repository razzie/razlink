package main

import (
	"fmt"
	"time"
)

// Log ...
type Log struct {
	Time     time.Time
	IP       string
	Location Location
}

// NewLog ...
func NewLog(ip string) Log {
	l := Log{Time: time.Now(), IP: ip}

	loc, _ := GetLocation(ip)
	if loc != nil {
		l.Location = *loc
	}

	return l
}

func (l Log) String() string {
	return fmt.Sprintf("%s - %s (%s)", l.Time, l.IP, l.Location.String())
}
