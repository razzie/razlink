package main

// Entry ...
type Entry struct {
	ID           string
	URL          string
	Method       ServeMethod
	Salt         string
	PasswordHash string
}

// NewEntry ...
func NewEntry(url, password string, method ServeMethod) *Entry {
	id := NewID()
	salt := id

	return &Entry{
		ID:           id,
		URL:          url,
		Method:       method,
		Salt:         salt,
		PasswordHash: Hash(salt + password),
	}
}

// MatchPassword ...
func (entry *Entry) MatchPassword(password string) bool {
	return entry.PasswordHash == Hash(entry.Salt+password)
}
