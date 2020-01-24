package razlink

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
	Hostnames   []string `json:"Addresses"`
	CountryName string
	RegionName  string
	City        string
	OS          string
	Browser     string
	Referer     string
}

// NewLog ...
func NewLog(r *http.Request) Log {
	ip := r.Header.Get("X-REAL-IP")
	if len(ip) == 0 {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	l := Log{
		Time:    time.Now(),
		IP:      ip,
		Referer: r.Referer(),
	}

	l.Hostnames, _ = net.LookupAddr(ip)
	if len(l.Hostnames) > 5 { // this would normally never happen, but let's make sure anyway
		l.Hostnames = l.Hostnames[:5]
	}

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
	return fmt.Sprintf("%s - %s - %s - %s/%s/%s - %s - %s",
		l.Time.Format(time.RFC3339),
		l.IP,
		strings.Join(l.Hostnames, ", "),
		l.CountryName,
		l.RegionName,
		l.City,
		l.OS,
		l.Browser)
}
