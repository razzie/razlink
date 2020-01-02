package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mssola/user_agent"
)

// Log ...
type Log struct {
	Time        time.Time
	IP          string
	Addresses   []string
	CountryName string
	RegionName  string
	City        string
	OS          string
	Browser     string
	Path        string `json:"-"`
}

// NewLog ...
func NewLog(r *http.Request) Log {
	ip := r.Header.Get("X-REAL-IP")
	l := Log{Time: time.Now(), IP: ip, Path: r.URL.Path}

	l.Addresses, _ = net.LookupAddr(ip)

	loc, _ := GetLocation(ip)
	if loc != nil {
		l.CountryName = loc.CountryName
		l.RegionName = loc.RegionName
		l.City = loc.City
	}

	ua := user_agent.New(r.UserAgent())
	browser, ver := ua.Browser()
	l.OS = ua.OS()
	l.Browser = fmt.Sprintf("%s %s", browser, ver)

	return l
}

func (l Log) String() string {
	return fmt.Sprintf("%s - %s - %s - %s/%s/%s - %s - %s - path: %s",
		l.Time.Format(time.RFC3339),
		l.IP,
		strings.Join(l.Addresses, ", "),
		l.CountryName,
		l.RegionName,
		l.City,
		l.OS,
		l.Browser,
		l.Path)
}
