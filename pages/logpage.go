package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/razzie/razlink"
)

var logAuthPageT = `
<form method="post">
	Log password:<br />
	<input type="password" name="password" /><br />
	<br />
	<input type="submit" value="Submit" />
</form>
`

var logPageT = `
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
	{{$ID := .ID}}
	{{range .Pages}}
		<a href="/logs/{{$ID}}/{{.}}">{{.}}</a> |
	{{end}}
	<a href="/logs/{{$ID}}/clear">clear</a>
{{else}}
	<strong>No logs yet!</strong>
{{end}}
`

func handleLogAuthPage(db *razlink.DB, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView("Not found", http.StatusNotFound)
	}

	if r.Method == "POST" {
		r.ParseForm()
		pw := r.FormValue("password")

		if !e.MatchPassword(pw) {
			return razlink.ErrorView("Wrong password", http.StatusUnauthorized)
		}

		cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
		return razlink.CookieAndRedirectView(r, cookie, "/logs/"+id)
	}

	return view(nil)
}

func handleLogPage(db *razlink.DB, logsPerPage int, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView("Not found", http.StatusNotFound)
	}

	cookie, _ := r.Cookie(id)
	if cookie == nil || cookie.Value != e.PasswordHash {
		return razlink.RedirectView(r, "/logs-auth/"+id)
	}

	if len(r.URL.Path) < 6+len(id)+1 { // /logs/ID/
		return razlink.RedirectView(r, "/logs/"+id+"/1")
	}

	actionOrPage := r.URL.Path[6+len(id)+1:]

	if actionOrPage == "clear" {
		db.DeleteLogs(id)
		return razlink.RedirectView(r, "/logs/"+id+"/1")
	}

	var data struct {
		ID    string
		Logs  []razlink.Log
		Pages []int
	}

	data.ID = id

	// pages
	logsCount, _ := db.GetLogsCount(id)
	pageCount := (logsCount / logsPerPage) + 1
	if logsCount > 0 && logsCount%logsPerPage == 0 {
		pageCount--
	}
	data.Pages = make([]int, pageCount)
	for i := range data.Pages {
		data.Pages[i] = i + 1
	}

	// logs
	page, _ := strconv.Atoi(actionOrPage)
	if page < 1 {
		page = 1
	} else if page > pageCount {
		return razlink.RedirectView(r, fmt.Sprintf("/logs/%s/%d", id, pageCount))
	}
	data.Logs, _ = db.GetLogs(id, (page-1)*logsPerPage, (page*logsPerPage)-1)

	return view(&data)
}

// GetLogPages ...
func GetLogPages(db *razlink.DB, logsPerPage int) []*razlink.Page {
	return []*razlink.Page{
		&razlink.Page{
			Path:            "/logs/",
			Title:           "Logs",
			ContentTemplate: logPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogPage(db, logsPerPage, r, view)
			},
		},
		&razlink.Page{
			Path:            "/logs-auth/",
			Title:           "Logs authentication",
			ContentTemplate: logAuthPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogAuthPage(db, r, view)
			},
		},
	}
}
