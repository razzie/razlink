package main

import (
	"fmt"
	"html/template"
	"net/http"
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
			<br />
			Access logs:<br />
			<a href="http://{{.Hostname}}/logs/{{.ID}}">{{.Hostname}}/logs/{{.ID}}</a>
		</div>
	</div>
</div>
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

		e, err := db.InsertEntry(url, pw, method)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.InsertLog(e.ID, r.Header.Get("X-REAL-IP"))

		http.Redirect(w, r, "/add/"+e.ID, http.StatusSeeOther)
	})

	mux.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		id, _ := getIDFromRequest(r)
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
