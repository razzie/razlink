package main

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
	id := NewID()
	salt := id

	return &Entry{
		ID:           id,
		URL:          url,
		Proxy:        proxy,
		Salt:         salt,
		PasswordHash: Hash(salt + password),
	}
}

// MatchPassword ...
func (entry *Entry) MatchPassword(password string) bool {
	return entry.PasswordHash == Hash(entry.Salt+password)
}

// Log ...
func (entry *Entry) Log(ip string) {
	l := NewLog(ip)
	entry.Logs = append(entry.Logs, l)
}
