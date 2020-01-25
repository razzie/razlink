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
type ViewFunc func(data interface{}) PageView

// PageHandler handles the page's requests
// If there was no error, the handler should call viewFunc and return the resulted PageView
type PageHandler func(r *http.Request, viewFunc ViewFunc) PageView

// BindLayout creates a page renderer function that uses Razlink layout
func (page *Page) BindLayout() (func(http.ResponseWriter, *http.Request), error) {
	renderer, err := BindLayout(page.ContentTemplate)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		viewFunc := func(data interface{}) PageView {
			return func(w http.ResponseWriter) {
				renderer(w, r, page.Title, data)
			}
		}

		var view PageView
		if page.Handler == nil {
			view = viewFunc(nil)
		} else {
			view = page.Handler(r, viewFunc)
		}

		view(w)
	}, nil
}
