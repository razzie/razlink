package razlink

import (
	"net/http"
)

// Page ...
type Page struct {
	Path            string
	Title           string
	ContentTemplate string
	Stylesheets     []string
	Scripts         []string
	Meta            map[string]string
	Handler         PageHandler
}

// ViewFunc is a function that produces a View using the input data
type ViewFunc func(data interface{}, title *string) *View

// PageHandler handles the page's requests
// If there was no error, the handler should call viewFunc and return the resulted PageView
type PageHandler func(r *http.Request, viewFunc ViewFunc) *View

// BindLayout creates a http.HandlerFunc that uses Razlink layout
func (page *Page) BindLayout() (http.HandlerFunc, error) {
	layoutRenderer, err := BindLayout(page.ContentTemplate, page.Stylesheets, page.Scripts, page.Meta)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		viewFunc := func(data interface{}, title *string) *View {
			renderer := func(w http.ResponseWriter) {
				if title != nil {
					layoutRenderer(w, r, *title, data)
				} else {
					layoutRenderer(w, r, page.Title, data)
				}
			}
			return &View{
				StatusCode: http.StatusOK,
				Data:       data,
				renderer:   renderer,
			}
		}

		var view *View
		if page.Handler == nil {
			view = viewFunc(nil, nil)
		} else {
			view = page.Handler(r, viewFunc)
		}

		view.Render(w)
	}, nil
}
