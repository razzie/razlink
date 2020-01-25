package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/razzie/razlink"
)

var logAuthPageT = `
{{if .}}
<strong style="color: red">{{.}}</strong><br /><br />
{{end}}
<form method="post">
	<input type="password" name="password" placeholder="Password" /><br />
	<button>Enter</button>
</form>
`

var logPageT = `
{{if .Logs}}
	<table>
		<tr>
			<td>Time</td>
			<td>IP</td>
			<td>Hostnames</td>
			<td>Country</td>
			<td>Region</td>
			<td>City</td>
			<td>OS</td>
			<td>Browser</td>
			<td>Referer</td>
		</tr>
		{{range .Logs}}
		<tr>
			<td>{{.Time.Format "Mon, 02 Jan 2006 15:04:05 MST"}}</td>
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
	<form method="get">
		{{$ID := .ID}}
		{{range .Pages}}
			<button formaction="/logs/{{$ID}}/{{.}}">{{.}}</button>
		{{end}}
		<button formaction="/logs-clear/{{$ID}}">clear</button>
	</form>
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
			return view("Wrong password")
		}

		cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
		return razlink.CookieAndRedirectView(r, cookie, "/logs/"+id)
	}

	return view(nil)
}

func handleLogPage(db *razlink.DB, logsPerPage int, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, trailing := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView("Not found", http.StatusNotFound)
	}

	cookie, _ := r.Cookie(id)
	if cookie == nil || cookie.Value != e.PasswordHash {
		return razlink.RedirectView(r, "/logs-auth/"+id)
	}

	pageCount := getLogPageCount(db, id, logsPerPage)

	page, _ := strconv.Atoi(trailing)
	if page < 1 {
		return razlink.RedirectView(r, fmt.Sprintf("/logs/%s/1", id))
	} else if page > pageCount {
		return razlink.RedirectView(r, fmt.Sprintf("/logs/%s/%d", id, pageCount))
	}

	logs := getLogs(db, id, page, logsPerPage)

	return view(newLogPageData(id, logs, pageCount))
}

func handleLogClear(db *razlink.DB, r *http.Request) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView("Not found", http.StatusNotFound)
	}

	cookie, _ := r.Cookie(id)
	if cookie == nil || cookie.Value != e.PasswordHash {
		return razlink.RedirectView(r, "/logs-auth/"+id)
	}

	db.DeleteLogs(id)
	return razlink.RedirectView(r, "/logs/"+id+"/1")
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
		&razlink.Page{
			Path: "/logs-clear/",
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogClear(db, r)
			},
		},
	}
}

func getLogPageCount(db *razlink.DB, id string, logsPerPage int) int {
	logsCount, _ := db.GetLogsCount(id)
	pageCount := (logsCount / logsPerPage) + 1
	if logsCount > 0 && logsCount%logsPerPage == 0 {
		pageCount--
	}
	return pageCount
}

func getLogs(db *razlink.DB, id string, page, logsPerPage int) []razlink.Log {
	logs, _ := db.GetLogs(id, (page-1)*logsPerPage, (page*logsPerPage)-1)
	return logs
}

func newLogPageData(id string, logs []razlink.Log, pageCount int) interface{} {
	pages := make([]int, pageCount)
	for i := range pages {
		pages[i] = i + 1
	}

	return &struct {
		ID    string
		Logs  []razlink.Log
		Pages []int
	}{
		ID:    id,
		Logs:  logs,
		Pages: pages,
	}
}
