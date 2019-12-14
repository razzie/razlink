package main

// DB ...
type DB struct {
	entries map[string]*Entry
	logs    map[string][]Log
}

// NewDB returns a new DB
func NewDB() *DB {
	return &DB{
		entries: make(map[string]*Entry),
		logs:    make(map[string][]Log),
	}
}

// InsertEntry inserts a new entry to the database
func (db *DB) InsertEntry(url, password string, proxy bool) *Entry {
	e := NewEntry(url, password, proxy)
	db.entries[e.ID] = e
	return e
}

// GetEntry returns the entry with the given ID
func (db *DB) GetEntry(id string) *Entry {
	e, _ := db.entries[id]
	return e
}

// DeleteEntry deleted the entry with the given ID
func (db *DB) DeleteEntry(id string) {
	delete(db.entries, id)
}

// InsertLog inserts a new log
func (db *DB) InsertLog(entryID, ip string) {
	l := NewLog(ip)
	logs, _ := db.logs[entryID]
	db.logs[entryID] = append(logs, l)
}

// GetLogs returns the logs that belong to an entry
func (db *DB) GetLogs(entryID string) []Log {
	logs, _ := db.logs[entryID]
	return logs
}

// DeleteLogs deleted all logs that belong to an entry
func (db *DB) DeleteLogs(entryID string) {
	delete(db.logs, entryID)
}
