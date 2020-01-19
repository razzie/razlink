package pages

import (
	"fmt"
	"io"
	"net/http"

	"github.com/razzie/razlink"
)

func embed(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<iframe src="%s" style="position:fixed; top:0; left:0; bottom:0; right:0; width:100%%; height:100%%; border:none; margin:0; padding:0; overflow:hidden; z-index:999999;"></iframe>`, url)
}

func handleViewPage(db *razlink.DB, w http.ResponseWriter, r *http.Request) interface{} {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return nil
	}

	defer db.InsertLog(id, r)

	switch e.Method {
	case razlink.Proxy:
		req, _ := http.NewRequest("GET", e.URL, nil)
		resp, err := http.DefaultClient.Do(req.WithContext(r.Context()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		defer resp.Body.Close()

		// Success is indicated with 2xx status codes:
		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			http.Error(w, resp.Status, resp.StatusCode)
			return nil
		}

		// in case the served content is not a file anymore
		if razlink.HasContentType(resp.Header, "text/html") {
			if len(resp.Header.Get("X-Frame-Options")) > 0 {
				e.Method = razlink.Redirect
				http.Redirect(w, r, e.URL, http.StatusSeeOther)
			} else {
				e.Method = razlink.Embed
				embed(w, e.URL)
			}

			db.SetEntry(id, e)
			return nil
		}

		for k, v := range resp.Header {
			w.Header().Set(k, v[0])
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

	case razlink.Embed:
		embed(w, e.URL)

	case razlink.Redirect:
		http.Redirect(w, r, e.URL, http.StatusSeeOther)

	case razlink.Track:
		razlink.WritePixel(w)

	default:
		http.Error(w, "Invalid serve method", http.StatusInternalServerError)
	}

	return nil
}

// GetViewPage ...
func GetViewPage(db *razlink.DB) *razlink.Page {
	return &razlink.Page{
		Path: "/x/",
		Handler: func(w http.ResponseWriter, r *http.Request) interface{} {
			return handleViewPage(db, w, r)
		},
	}
}
