package main

import (
	"encoding/json"
	"fmt"
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

	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	success, err := db.client.SetNX(e.ID, string(data), db.ExpirationTime).Result()
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, fmt.Errorf("duplicate ID: %s", e.ID)
	}

	return e, nil
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

	db.client.Expire(id, db.ExpirationTime) // reset expiration
	db.client.Expire(id+"-log", db.ExpirationTime)

	return &e, nil
}

// DeleteEntry deleted the entry with the given ID
func (db *DB) DeleteEntry(id string) error {
	defer db.DeleteLogs(id)
	return db.client.Del(id).Err()
}

// InsertLog inserts a new log
func (db *DB) InsertLog(entryID, ip string) error {
	l := NewLog(ip)

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
func (db *DB) GetLogs(entryID string, first, last uint) ([]Log, error) {
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
func (db *DB) GetLogsCount(entryID string) (uint, error) {
	value, err := db.client.LLen(entryID + "-log").Result()
	return uint(value), err
}

// DeleteLogs deleted all logs that belong to an entry
func (db *DB) DeleteLogs(entryID string) error {
	return db.client.Del(entryID + "-log").Err()
}
