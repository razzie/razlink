package main

import (
	"fmt"
	"time"
)

// Log ...
type Log struct {
	Time        time.Time
	IP          string
	CountryName string
	RegionName  string
	City        string
}

// NewLog ...
func NewLog(ip string) Log {
	l := Log{Time: time.Now(), IP: ip}

	loc, _ := GetLocation(ip)
	if loc != nil {
		l.CountryName = loc.CountryName
		l.RegionName = loc.RegionName
		l.City = loc.City
	}

	return l
}

func (l Log) String() string {
	return fmt.Sprintf("%s - %s (%s/%s/%s)",
		l.Time.Format(time.RFC3339),
		l.IP,
		l.CountryName,
		l.RegionName,
		l.City)
}
