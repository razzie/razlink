package pages

import (
	"net/http"

	"github.com/razzie/razlink"
)

func handleViewPage(db *razlink.DB, r *http.Request) *razlink.View {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView(r, "Not found", http.StatusNotFound)
	}

	defer db.InsertLog(id, r)

	switch e.Method {
	case razlink.Proxy:
		req, _ := http.NewRequest("GET", e.URL, nil)
		resp, err := http.DefaultClient.Do(req.WithContext(r.Context()))
		if err != nil {
			return razlink.ErrorView(r, err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		// Success is indicated with 2xx status codes:
		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			return razlink.ErrorView(r, resp.Status, resp.StatusCode)
		}

		// in case the served content is not a file anymore
		if razlink.HasContentType(resp.Header, "text/html") {
			defer db.SetEntry(id, e)

			if len(resp.Header.Get("X-Frame-Options")) > 0 {
				e.Method = razlink.Redirect
				return razlink.RedirectView(r, e.URL)
			}

			e.Method = razlink.Embed
			return razlink.EmbedView(e.URL)
		}

		return razlink.CopyView(resp)

	case razlink.Embed:
		return razlink.EmbedView(e.URL)

	case razlink.Redirect:
		return razlink.RedirectView(r, e.URL)

	case razlink.Track:
		return razlink.HandlerView(r, razlink.WritePixel)

	default:
		return razlink.ErrorView(r, "Invalid serve method", http.StatusInternalServerError)
	}
}

// GetViewPage ...
func GetViewPage(db *razlink.DB) *razlink.Page {
	return &razlink.Page{
		Path: "/x/",
		Handler: func(r *http.Request, view razlink.ViewFunc) *razlink.View {
			return handleViewPage(db, r)
		},
	}
}
