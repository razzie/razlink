package main

// Entry ...
type Entry struct {
	ID           string
	URL          string
	Method       ServeMethod
	Salt         string
	PasswordHash string
	Permanent    bool
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

// NewPermanentEntry ...
func NewPermanentEntry(ID, url, password string, method ServeMethod) *Entry {
	e := NewEntry(url, password, method)
	e.ID = ID
	e.Permanent = true
	return e
}

// MatchPassword ...
func (entry *Entry) MatchPassword(password string) bool {
	return entry.PasswordHash == Hash(entry.Salt+password)
}
