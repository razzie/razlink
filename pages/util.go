package pages

import (
	"net/http"
	"strings"
)

func getIDFromRequest(r *http.Request) (string, bool) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return "", false
	}
	return parts[2], true
}
