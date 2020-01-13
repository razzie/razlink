package pages

import (
	"html/template"
	"io"
	"net/http"

	"github.com/razzie/razlink"
)

var embedPage = `
<iframe src="{{.}}" style="position:fixed; top:0; left:0; bottom:0; right:0; width:100%; height:100%; border:none; margin:0; padding:0; overflow:hidden; z-index:999999;"></iframe>
`

func installViewPage(db *razlink.DB, mux *http.ServeMux) {
	embedPageT, err := template.New("").Parse(embedPage)
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/x/", func(w http.ResponseWriter, r *http.Request) {
		id, _ := getIDFromRequest(r)
		e, _ := db.GetEntry(id)
		if e == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		defer db.InsertLog(id, r)

		switch e.Method {
		case razlink.Proxy:
			req, _ := http.NewRequest("GET", e.URL, nil)
			resp, err := http.DefaultClient.Do(req.WithContext(r.Context()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				http.Error(w, resp.Status, resp.StatusCode)
				return
			}

			for k, v := range resp.Header {
				w.Header().Set(k, v[0])
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)

		case razlink.Embed:
			embedPageT.Execute(w, e.URL)

		case razlink.Redirect:
			http.Redirect(w, r, e.URL, http.StatusSeeOther)

		default:
			http.Error(w, "Invalid serve method", http.StatusInternalServerError)
		}
	})
}
