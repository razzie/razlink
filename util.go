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

// WriteFavicon writes Razlink favicon to a http.ResponseWriter
func WriteFavicon(w http.ResponseWriter) {
	favicon := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M18 16.08c-.76 0-1.44.3-1.96.77L8.91 12.7c.05-.23.09-.46.09-.7s-.04-.47-.09-.7l7.05-4.11c.54.5 1.25.81 2.04.81 1.66 0 3-1.34 3-3s-1.34-3-3-3-3 1.34-3 3c0 .24.04.47.09.7L8.04 9.81C7.5 9.31 6.79 9 6 9c-1.66 0-3 1.34-3 3s1.34 3 3 3c.79 0 1.5-.31 2.04-.81l7.12 4.16c-.05.21-.08.43-.08.65 0 1.61 1.31 2.92 2.92 2.92 1.61 0 2.92-1.31 2.92-2.92s-1.31-2.92-2.92-2.92z"></path></svg>`)
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Content-Length", strconv.Itoa(len(favicon)))
	_, _ = w.Write(favicon)
}

// GetBase returns the base target for relative URLs
func GetBase(r *http.Request) string {
	slashes := strings.Count(r.URL.Path, "/")
	if slashes > 1 {
		return strings.Repeat("../", slashes-1)
	}
	return "/"
}
