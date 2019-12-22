package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var logsPasswordPage = `
<div style="display: flex; align-items: center; justify-content: center">
	<div style="border: 1px solid black; padding: 1rem; display: inline-flex">
		<form method="post">
			Log password:<br />
			<input type="password" name="password" /><br />
			<br />
			<input type="submit" value="Submit" />
		</form>
	</div>
</div>
`

var logsPage = `
<div style="display: flex; align-items: center; justify-content: center">
	<div style="border: 1px solid black; padding: 1rem; display: inline-flex">
		<ul>
			{{range .}}
			<li>{{.}}</li>
			{{end}}
		</ul>
	</div>
</div>
`

func installLogPage(db *DB, mux *http.ServeMux) {
	logsPageT, err := template.New("").Parse(logsPage)
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Method != "POST" {
			fmt.Fprint(w, logsPasswordPage)
			return
		}

		r.ParseForm()
		pw := r.FormValue("password")

		id, _ := getIDFromRequest(r)
		e, _ := db.GetEntry(id)
		if e == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if !e.MatchPassword(pw) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logs, _ := db.GetLogs(id, 0)
		logsPageT.Execute(w, logs)
	})
}
