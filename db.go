package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

// DB ...
type DB struct {
	client *redis.Client
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
		client: client,
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

	err = db.client.Set(e.ID, string(data), 30*24*time.Hour).Err()
	if err != nil {
		return nil, err
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

	return &e, nil
}

// DeleteEntry deleted the entry with the given ID
func (db *DB) DeleteEntry(id string) error {
	return db.client.Del(id).Err()
}

// InsertLog inserts a new log
func (db *DB) InsertLog(entryID, ip string) error {
	l := NewLog(ip)

	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return db.client.LPush(entryID+"-log", string(data)).Err()
}

// GetLogs returns the logs that belong to an entry
func (db *DB) GetLogs(entryID string, page uint) ([]Log, error) {
	values, err := db.client.LRange(entryID+"-log", int64(page)*100, int64(page+1)*100).Result()
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

// DeleteLogs deleted all logs that belong to an entry
func (db *DB) DeleteLogs(entryID string) error {
	return db.client.Del(entryID + "-log").Err()
}
