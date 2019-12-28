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
		<table>
			<tr>
				<th>Time</th>
				<th>IP</th>
				<th>Addresses</th>
				<th>Country</th>
				<th>Region</th>
				<th>City</th>
			</tr>
			{{range .}}
			<tr>
				<td>{{.Time}}</td>
				<td>{{.IP}}</td>
				<td>
				{{range .Addresses}}
					{{.}}<br />
				{{end}}
				</td>
				<td>{{.CountryName}}</td>
				<td>{{.RegionName}}</td>
				<td>{{.City}}</td>
			</tr>
			{{end}}
		</table>
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

		id, _ := getIDFromRequest(r)
		e, _ := db.GetEntry(id)
		if e == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if r.Method == "POST" {
			r.ParseForm()
			pw := r.FormValue("password")

			if !e.MatchPassword(pw) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, &http.Cookie{Name: e.ID, Value: e.PasswordHash})
		} else {
			cookie, _ := r.Cookie(e.ID)
			if cookie == nil || cookie.Value != e.PasswordHash {
				fmt.Fprint(w, logsPasswordPage)
				return
			}
		}

		logs, _ := db.GetLogs(id, 0, 100)
		logsPageT.Execute(w, logs)
	})
}
