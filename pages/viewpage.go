package pages

import (
	"net/http"

	"github.com/razzie/razlink"
)

func handleViewPage(db *razlink.DB, pr *razlink.PageRequest) *razlink.View {
	id, _ := getIDFromRequest(pr)
	e, _ := db.GetEntry(id)
	if e == nil {
		return pr.ErrorView("Not found", http.StatusNotFound)
	}

	r := pr.Request
	defer db.InsertLog(id, r)

	switch e.Method {
	case razlink.Proxy:
		req, _ := http.NewRequest("GET", e.URL, nil)
		resp, err := http.DefaultClient.Do(req.WithContext(r.Context()))
		if err != nil {
			return pr.ErrorView(err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		// Success is indicated with 2xx status codes:
		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			return pr.ErrorView(resp.Status, resp.StatusCode)
		}

		// in case the served content is not a file anymore
		if razlink.HasContentType(resp.Header, "text/html") {
			defer db.SetEntry(id, e)

			if len(resp.Header.Get("X-Frame-Options")) > 0 {
				e.Method = razlink.Redirect
				return pr.RedirectView(e.URL)
			}

			e.Method = razlink.Embed
			return pr.EmbedView(e.URL)
		}

		return pr.CopyView(resp)

	case razlink.Embed:
		return pr.EmbedView(e.URL)

	case razlink.Redirect:
		return pr.RedirectView(e.URL)

	case razlink.Track:
		return pr.HandlerView(razlink.WritePixel)

	default:
		return pr.ErrorView("Invalid serve method", http.StatusInternalServerError)
	}
}

// GetViewPage ...
func GetViewPage(db *razlink.DB) *razlink.Page {
	return &razlink.Page{
		Path: "/x/",
		Handler: func(pr *razlink.PageRequest) *razlink.View {
			return handleViewPage(db, pr)
		},
	}
}
