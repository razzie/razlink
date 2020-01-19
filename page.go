package razlink

import (
	"net/http"
)

// Page ...
type Page struct {
	Path            string
	Title           string
	ContentTemplate string
	Handler         PageHandler
}

// PageHandler handles the page's requests
// Returns data to be used by the page template or nil if there is no need to render the template
type PageHandler func(w http.ResponseWriter, r *http.Request) interface{}
