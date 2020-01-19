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

// PageView is a callback function used to render the page
type PageView func(w http.ResponseWriter)

// ViewFunc is a function that produces a PageView using the input data
type ViewFunc func(interface{}) PageView

// PageHandler handles the page's requests
// If there was no error, the handler should call viewFunc and return the resulted PageView
type PageHandler func(r *http.Request, viewFunc ViewFunc) PageView
