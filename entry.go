package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
)

// Entry ...
type Entry struct {
	ID           string
	URL          string
	Proxy        bool
	Salt         string
	PasswordHash string
	Logs         []Log
}

// NewEntry ...
func NewEntry(url, password string, proxy bool) *Entry {
	id := strconv.FormatInt(time.Now().UnixNano(), 16)
	salt := id

	return &Entry{
		ID:           id,
		URL:          url,
		Proxy:        proxy,
		Salt:         salt,
		PasswordHash: hash(salt + password),
	}
}

// MatchPassword ...
func (entry *Entry) MatchPassword(password string) bool {
	return entry.PasswordHash == hash(entry.Salt+password)
}

// Log ...
func (entry *Entry) Log(ip string) {
	l := NewLog(ip)
	entry.Logs = append(entry.Logs, l)
}

func hash(s string) string {
	algorithm := sha1.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}
