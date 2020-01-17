package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

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
			{{if .Track}}
			Embed this in your website:<br />
			&lt;img src="<a href="http://{{.Hostname}}/x/{{.ID}}">http://{{.Hostname}}/x/{{.ID}}</a>" width="1" height="1" /&gt;<br />
			{{else}}
			Access the target URL:<br />
			<a href="http://{{.Hostname}}/x/{{.ID}}">{{.Hostname}}/x/{{.ID}}</a><br />
			{{if .Decoy}}
			<a href="http://{{.Hostname}}/x/{{.ID}}/{{.Decoy}}">{{.Hostname}}/x/{{.ID}}/{{.Decoy}}</a><br />
			{{end}}
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
		method, err := razlink.GetServeMethodForURL(r.Context(), url, time.Second)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := razlink.NewEntry(url, pw, method)
		id, err := db.InsertEntry(nil, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.InsertLog(id, r)

		http.SetCookie(w, &http.Cookie{Name: id, Value: e.PasswordHash})
		http.Redirect(w, r, "/add/"+id, http.StatusSeeOther)
	})

	mux.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		id, _ := getIDFromRequest(r)

		e, _ := db.GetEntry(id)
		if e == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		decoy := filepath.Base(e.URL)
		if len(decoy) < 2 {
			decoy = ""
		}

		view := struct {
			Hostname string
			ID       string
			Decoy    string
			Track    bool
		}{
			Hostname: hostname,
			ID:       id,
			Decoy:    decoy,
			Track:    e.Method == razlink.Track,
		}

		addResultPageT.Execute(w, view)
	})
}
