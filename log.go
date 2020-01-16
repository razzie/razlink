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
	Hostnames   []string
	CountryName string
	RegionName  string
	City        string
	OS          string
	Browser     string
	Referer     string
	Path        string `json:"-"`
}

// NewLog ...
func NewLog(r *http.Request) Log {
	ip := r.Header.Get("X-REAL-IP")
	l := Log{
		Time:    time.Now(),
		IP:      ip,
		Referer: r.Referer(),
		Path:    r.URL.Path,
	}

	l.Hostnames, _ = net.LookupAddr(ip)

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
		strings.Join(l.Hostnames, ", "),
		l.CountryName,
		l.RegionName,
		l.City,
		l.OS,
		l.Browser,
		l.Path)
}
