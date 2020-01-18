package razlink

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var instance string
var privateIPBlocks []*net.IPNet

func init() {
	// set up a unique ID for this instance (to be prepended to entry IDs)
	i := uint16(time.Now().UnixNano())
	instance = strconv.FormatInt(int64(i), 36)

	// initializing private IP address spaces
	// https://stackoverflow.com/a/50825191
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
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

// IsPrivateIP returns true if the given IP address belongs to private network space
func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// HasContentType determines whether the request `content-type` includes a
// server-acceptable mime-type
func HasContentType(header http.Header, mimetype string) bool {
	contentType := header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}

// WritePixel writes a transparent pixel to a http.ResponseWriter
func WritePixel(w http.ResponseWriter) {
	pixel := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="
	bytes, _ := base64.StdEncoding.DecodeString(pixel)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	_, _ = w.Write(bytes)
}
