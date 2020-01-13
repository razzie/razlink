package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

// DB ...
type DB struct {
	ExpirationTime time.Duration
	client         *redis.Client
}

// NewDB returns a new DB
func NewDB(addr, password string, db int) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	err := client.Ping().Err()
	if err != nil {
		client.Close()
		return nil, err
	}

	return &DB{
		ExpirationTime: 30 * 24 * time.Hour, // ~1 month
		client:         client,
	}, nil
}

// Close closes the connection to the database
func (db *DB) Close() error {
	return db.client.Close()
}

// InsertEntry inserts a new entry to the database
func (db *DB) InsertEntry(url, password string, method ServeMethod) (*Entry, error) {
	e := NewEntry(url, password, method)
	return e, db.insertEntry(e)
}

// InsertPermanentEntry inserts a new entry to the database
func (db *DB) InsertPermanentEntry(ID, url, password string, method ServeMethod) (*Entry, error) {
	e := NewPermanentEntry(ID, url, password, method)
	return e, db.insertEntry(e)
}

func (db *DB) insertEntry(e *Entry) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	expiration := db.ExpirationTime
	if e.Permanent {
		expiration = 0
	}

	success, err := db.client.SetNX(e.ID, string(data), expiration).Result()
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("duplicate ID: %s", e.ID)
	}

	return nil
}

// GetEntry returns the entry with the given ID
func (db *DB) GetEntry(id string) (*Entry, error) {
	data, err := db.client.Get(id).Result()
	if err != nil {
		return nil, err
	}

	var e Entry
	err = json.Unmarshal([]byte(data), &e)
	if err != nil {
		return nil, err
	}

	// reset expiration
	if !e.Permanent {
		db.client.Expire(id, db.ExpirationTime)
	}
	db.client.Expire(id+"-log", db.ExpirationTime)

	return &e, nil
}

// GetEntries returns the list of entries with IDs matching the given pattern
func (db *DB) GetEntries(pattern string) ([]*Entry, error) {
	keys, err := db.client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}

	var entries []*Entry
	for _, key := range keys {
		if strings.HasSuffix(key, "-log") {
			continue
		}

		e, err := db.GetEntry(key)
		if err != nil {
			continue
		}

		entries = append(entries, e)
	}

	return entries, nil
}

// DeleteEntry deleted the entry with the given ID
func (db *DB) DeleteEntry(id string) error {
	defer db.DeleteLogs(id)
	return db.client.Del(id).Err()
}

// InsertLog inserts a new log
func (db *DB) InsertLog(entryID string, r *http.Request) error {
	l := NewLog(r)

	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	len, err := db.client.LPush(entryID+"-log", string(data)).Result()
	if err != nil {
		return err
	}

	if len == 1 {
		db.client.Expire(entryID+"-log", db.ExpirationTime)
	}

	return nil
}

// GetLogs returns the Nth page of logs that belong to an entry (pages are 0 based)
func (db *DB) GetLogs(entryID string, first, last int) ([]Log, error) {
	values, err := db.client.LRange(entryID+"-log", int64(first), int64(last)).Result()
	if err != nil {
		return nil, err
	}

	var logs []Log
	for _, data := range values {
		var l Log
		err := json.Unmarshal([]byte(data), &l)
		if err == nil {
			logs = append(logs, l)
		}
	}

	return logs, nil
}

// GetLogsCount returns the number of logs that belong to an entry
func (db *DB) GetLogsCount(entryID string) (int, error) {
	value, err := db.client.LLen(entryID + "-log").Result()
	return int(value), err
}

// DeleteLogs deleted all logs that belong to an entry
func (db *DB) DeleteLogs(entryID string) error {
	return db.client.Del(entryID + "-log").Err()
}
