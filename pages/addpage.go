package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/razzie/razlink"
)

var addPage = `
<div style="display: flex; align-items: center; justify-content: center">
	<div style="border: 1px solid black; padding: 1rem; display: inline-flex">
		<form method="post">
			URL:<br />
			<input type="text" name="url" /><br />
			<br />
			Log password:<br />
			<input type="password" name="password" /><br />
			<br />
			<input type="submit" value="Submit" />
		</form>
	</div>
</div>
`

var addResultPage = `
<div style="display: flex; align-items: center; justify-content: center">
	<div style="border: 1px solid black; padding: 1rem; display: inline-flex">
		<div>
			<strong>Bookmark this page!</strong><br />
			<br />
			Access the target URL:<br />
			<a href="http://{{.Hostname}}/x/{{.ID}}">{{.Hostname}}/x/{{.ID}}</a><br />
			{{if .Decoy}}
			<a href="http://{{.Hostname}}/x/{{.ID}}{{.Decoy}}">{{.Hostname}}/x/{{.ID}}{{.Decoy}}</a><br />
			{{end}}
			<br />
			Access logs:<br />
			<a href="http://{{.Hostname}}/logs/{{.ID}}">{{.Hostname}}/logs/{{.ID}}</a>
		</div>
	</div>
</div>
`

func installAddPage(db *razlink.DB, mux *http.ServeMux, hostname string) {
	addResultPageT, err := template.New("").Parse(addResultPage)
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, addPage)
			return
		}

		r.ParseForm()
		url := r.FormValue("url")
		pw := r.FormValue("password")
		method, err := razlink.GetServeMethodForURL(r.Context(), url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		decoy := filepath.Base(url)
		if len(decoy) < 2 {
			decoy = ""
		}

		e := razlink.NewEntry(url, pw, method)
		id, err := db.InsertEntry(nil, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.InsertLog(id, r)

		http.SetCookie(w, &http.Cookie{Name: id, Value: e.PasswordHash})
		http.Redirect(w, r, fmt.Sprintf("/add/%s/%s", id, decoy), http.StatusSeeOther)
	})

	mux.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		id, _ := getIDFromRequest(r)
		view := struct {
			Hostname string
			ID       string
			Decoy    string
		}{
			Hostname: hostname,
			ID:       id,
			Decoy:    r.URL.Path[5+len(id):],
		}
		addResultPageT.Execute(w, view)
	})
}