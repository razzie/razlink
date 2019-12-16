package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var logsPasswordPage = `
<form method="post">
	<input type="password" name="password" />
	<input type="submit" value="Submit" />
</form>
`

var logsPage = `
<ul>
	{{range .}}
	<li>{{.}}</li>
	{{end}}
</ul>
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

		id := r.URL.Path[6:]
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
