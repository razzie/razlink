package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
		{{if .Logs}}
		<table>
			<tr>
				<th>Time</th>
				<th>IP</th>
				<th>Addresses</th>
				<th>Country</th>
				<th>Region</th>
				<th>City</th>
			</tr>
			{{range .Logs}}
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
			<tr>
				<td colspan="6">
					{{range .Pages}}
					<a href="{{.}}">{{.}}</a> |
					{{end}}
					<a href="clear">clear</a>
				</td>
			</tr>
		</table>
		{{else}}
		<strong>No logs yet!</strong>
		{{end}}
	</div>
</div>
`

const logsPerPage = 20

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

		if len(r.URL.Path) < 6+len(id)+1 { // /logs/ID/
			http.Redirect(w, r, "/logs/"+e.ID+"/1", http.StatusSeeOther)
			return
		}

		actionOrPage := r.URL.Path[6+len(id)+1:]

		if actionOrPage == "clear" {
			db.DeleteLogs(e.ID)
			http.Redirect(w, r, "/logs/"+e.ID+"/1", http.StatusSeeOther)
			return
		}

		var view struct {
			Logs  []Log
			Pages []int
		}

		// pages
		logsCount, _ := db.GetLogsCount(id)
		pageCount := (logsCount / logsPerPage) + 1
		if logsCount > 0 && logsCount%logsPerPage == 0 {
			pageCount--
		}
		view.Pages = make([]int, pageCount)
		for i := range view.Pages {
			view.Pages[i] = i + 1
		}

		// logs
		page, _ := strconv.Atoi(actionOrPage)
		if page < 1 {
			page = 1
		} else if page > pageCount {
			http.Redirect(w, r, fmt.Sprintf("/logs/%s/%d", e.ID, pageCount), http.StatusSeeOther)
			return
		}
		view.Logs, _ = db.GetLogs(id, (page-1)*logsPerPage, (page*logsPerPage)-1)

		logsPageT.Execute(w, &view)
	})
}
