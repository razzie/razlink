package pages

import (
	"fmt"
	"net/http"

	"github.com/razzie/razlink"
)

func embed(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<iframe src="%s" style="position:fixed; top:0; left:0; bottom:0; right:0; width:100%%; height:100%%; border:none; margin:0; padding:0; overflow:hidden; z-index:999999;"></iframe>`, url)
}

func handleViewPage(db *razlink.DB, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView("Not found", http.StatusNotFound)
	}

	defer db.InsertLog(id, r)

	switch e.Method {
	case razlink.Proxy:
		req, _ := http.NewRequest("GET", e.URL, nil)
		resp, err := http.DefaultClient.Do(req.WithContext(r.Context()))
		if err != nil {
			return razlink.ErrorView(err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		// Success is indicated with 2xx status codes:
		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			return razlink.ErrorView(resp.Status, resp.StatusCode)
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
		return razlink.WritePixel

	default:
		return razlink.ErrorView("Invalid serve method", http.StatusInternalServerError)
	}
}

// GetViewPage ...
func GetViewPage(db *razlink.DB) *razlink.Page {
	return &razlink.Page{
		Path: "/x/",
		Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
			return handleViewPage(db, r, view)
		},
	}
}
