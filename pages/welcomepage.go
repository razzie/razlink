package pages

import (
	"net/http"

	"github.com/razzie/razlink"
)

var welcomePageT = `
<strong>Welcome to razlink!</strong><br />
<br />
This is a lightweight link proxy/redirect service with logging.<br />
<br />
You can <a href="create">add a new link</a> with a custom password to view the visits.<br />
There are 4 automatic operation modes based on the URL:<br />
<ul>
	<li>Proxy - for files (e.g. an image)</li>
	<li>Embed - for websites that allow embedding</li>
	<li>Redirect - for websites that forbid embedding</li>
	<li>Track - write <strong>.</strong> to URL field to enable transparent pixel track mode</li>
</ul><br />
Check out the source code at <a href="https://github.com/razzie/razlink" target="_blank">github.com/razzie/razlink</a>.<br />
`

// GetWelcomePage ...
func GetWelcomePage() *razlink.Page {
	return &razlink.Page{
		Path:            "/",
		Title:           "Welcome to razlink!",
		ContentTemplate: welcomePageT,
		Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
			if len(r.URL.Path) > 1 {
				return razlink.RedirectView(r, "/")
			}
			return view(nil)
		},
	}
}
