package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Log ...
type Log struct {
	Time        time.Time
	IP          string
	Addresses   []string
	CountryName string
	RegionName  string
	City        string
}

// NewLog ...
func NewLog(ip string) Log {
	l := Log{Time: time.Now(), IP: ip}

	l.Addresses, _ = net.LookupAddr(ip)

	loc, _ := GetLocation(ip)
	if loc != nil {
		l.CountryName = loc.CountryName
		l.RegionName = loc.RegionName
		l.City = loc.City
	}

	return l
}

func (l Log) String() string {
	return fmt.Sprintf("%s - %s - %s - %s / %s / %s",
		l.Time.Format(time.RFC3339),
		l.IP,
		strings.Join(l.Addresses, ", "),
		l.CountryName,
		l.RegionName,
		l.City)
}
