package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var addPage = `
<form method="post">
	URL: <input type="text" name="url" /><br />
	Log password: <input type="password" name="password" /><br />
	<input type="submit" value="Submit" />
</form>
`

var addResultPage = `
<a href="http://{{.Hostname}}/x/{{.ID}}">{{.Hostname}}/x/{{.ID}}</a><br />
<a href="http://{{.Hostname}}/logs/{{.ID}}">{{.Hostname}}/logs/{{.ID}}</a>
`

func installAddPage(db *DB, mux *http.ServeMux, hostname string) {
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
		method, err := GetServeMethodForURL(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := db.InsertEntry(url, pw, method)
		http.Redirect(w, r, "/add/"+e.ID, http.StatusSeeOther)
	})

	mux.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		id := r.URL.Path[5:]
		view := struct {
			Hostname string
			ID       string
		}{
			Hostname: hostname,
			ID:       id,
		}
		addResultPageT.Execute(w, view)
	})
}
