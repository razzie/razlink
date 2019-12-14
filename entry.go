package main

// Entry ...
type Entry struct {
	ID           string
	URL          string
	Proxy        bool
	Salt         string
	PasswordHash string
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
