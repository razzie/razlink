package pages

import (
	"net/http"
	"strings"
)

func getIDFromRequest(r *http.Request) (id string, trailing string) {
	parts := strings.SplitN(r.URL.Path, "/", 4) // example: /logs/<id>/<page>
	if len(parts) >= 3 {
		id = parts[2]
		if len(parts) >= 4 {
			trailing = parts[3]
		}
	}
	return
}
