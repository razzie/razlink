package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
)

var instance string

func init() {
	i := uint16(time.Now().UnixNano())
	instance = strconv.FormatInt(int64(i), 36)
}

// NewID returns a new (hopefully unique) ID for entries
func NewID() string {
	return instance + "-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}

// Hash returns the SHA1 hash of a string
func Hash(s string) string {
	algorithm := sha1.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}
