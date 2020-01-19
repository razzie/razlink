package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/razzie/razlink"
)

var logAuthPageT = `
{{define "page"}}
<form method="post">
	Log password:<br />
	<input type="password" name="password" /><br />
	<br />
	<input type="submit" value="Submit" />
</form>
{{end}}
`

var logPageT = `
{{define "page"}}
{{if .Logs}}
	<table style="border-spacing: 10px; border-collapse: separate">
		<tr>
			<td><strong>Time</strong></td>
			<td><strong>IP</strong></td>
			<td><strong>Hostnames</strong></td>
			<td><strong>Country</strong></td>
			<td><strong>Region</strong></td>
			<td><strong>City</strong></td>
			<td><strong>OS</strong></td>
			<td><strong>Browser</strong></td>
			<td><strong>Referer</strong></td>
		</tr>
		{{range .Logs}}
		<tr>
			<td>{{.Time}}</td>
			<td>{{.IP}}</td>
			<td>
				{{range .Hostnames}}
				{{.}}<br />
				{{end}}
			</td>
			<td>{{.CountryName}}</td>
			<td>{{.RegionName}}</td>
			<td>{{.City}}</td>
			<td>{{.OS}}</td>
			<td>{{.Browser}}</td>
			<td>{{.Referer}}</td>
		</tr>
		{{end}}
	</table>
	{{range .Pages}}
		<a href="{{.}}">{{.}}</a> |
	{{end}}
	<a href="clear">clear</a>
{{else}}
	<strong>No logs yet!</strong>
{{end}}
{{end}}
`

func handleLogAuthPage(db *razlink.DB, w http.ResponseWriter, r *http.Request) interface{} {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return nil
	}

	if r.Method == "POST" {
		r.ParseForm()
		pw := r.FormValue("password")

		if !e.MatchPassword(pw) {
			http.Error(w, "Wrong password", http.StatusUnauthorized)
			return nil
		}

		http.SetCookie(w, &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"})
		http.Redirect(w, r, "/logs/"+id, http.StatusSeeOther)
		return nil
	}

	return ""
}

func handleLogPage(db *razlink.DB, logsPerPage int, w http.ResponseWriter, r *http.Request) interface{} {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return nil
	}

	cookie, _ := r.Cookie(id)
	if cookie == nil || cookie.Value != e.PasswordHash {
		http.Redirect(w, r, "/logs-auth/"+id, http.StatusSeeOther)
		return nil
	}

	if len(r.URL.Path) < 6+len(id)+1 { // /logs/ID/
		http.Redirect(w, r, "/logs/"+id+"/1", http.StatusSeeOther)
		return nil
	}

	actionOrPage := r.URL.Path[6+len(id)+1:]

	if actionOrPage == "clear" {
		db.DeleteLogs(id)
		http.Redirect(w, r, "/logs/"+id+"/1", http.StatusSeeOther)
		return nil
	}

	var view struct {
		Logs  []razlink.Log
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
		http.Redirect(w, r, fmt.Sprintf("/logs/%s/%d", id, pageCount), http.StatusSeeOther)
		return nil
	}
	view.Logs, _ = db.GetLogs(id, (page-1)*logsPerPage, (page*logsPerPage)-1)

	return &view
}

// GetLogPages ...
func GetLogPages(db *razlink.DB, logsPerPage int) []*razlink.Page {
	return []*razlink.Page{
		&razlink.Page{
			Path:            "/logs/",
			Title:           "Logs",
			ContentTemplate: logPageT,
			Handler: func(w http.ResponseWriter, r *http.Request) interface{} {
				return handleLogPage(db, logsPerPage, w, r)
			},
		},
		&razlink.Page{
			Path:            "/logs-auth/",
			Title:           "Logs authentication",
			ContentTemplate: logAuthPageT,
			Handler: func(w http.ResponseWriter, r *http.Request) interface{} {
				return handleLogAuthPage(db, w, r)
			},
		},
	}
}
