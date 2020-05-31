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

var logChangePwPageT = `
{{if .}}
<strong style="color: red">{{.}}</strong><br /><br />
{{end}}
<form method="post">
	<input type="password" name="old_password" placeholder="Old password" /><br />
	<input type="password" name="password" placeholder="New password" /><br />
	<input type="password" name="confirm_password" placeholder="Confirm new password" /><br />
	<button>Save</button>
</form>
`

var logPageT = `
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
	{{else}}
	<tr>
		<td colspan="9" style="text-align: center"><strong>No logs yet!</strong></td>
	</tr>
	{{end}}
</table>
<form method="get">
	{{$ID := .ID}}
	{{range .Pages}}
		<button formaction="/logs/{{$ID}}/{{.}}">{{.}}</button>
	{{end}}
	<button formaction="/logs-clear/{{$ID}}/" onclick="return confirm('Are you sure?')">clear</button>
	<button formaction="/logs-change-password/{{$ID}}/">change password</button>
</form>
`

func handleLogAuthPage(db *razlink.DB, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView(r, "Not found", http.StatusNotFound)
	}

	if r.Method == "POST" {
		r.ParseForm()
		pw := r.FormValue("password")

		if !e.MatchPassword(pw) {
			return view("Wrong password", nil)
		}

		cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
		return razlink.CookieAndRedirectView(r, cookie, "/logs/"+id)
	}

	return view(nil, nil)
}

func handleLogPage(db *razlink.DB, logsPerPage int, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, trailing := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView(r, "Not found", http.StatusNotFound)
	}

	cookie, _ := r.Cookie(id)
	if cookie == nil || cookie.Value != e.PasswordHash {
		return razlink.RedirectView(r, "/logs-auth/"+id)
	}

	var first, last int
	var logs []razlink.Log
	title := "Logs of " + id
	pageCount := getLogPageCount(db, id, logsPerPage)

	if _, err := fmt.Sscanf(trailing, "%d..%d", &first, &last); err != nil {
		page, _ := strconv.Atoi(trailing)
		if page < 1 {
			return razlink.RedirectView(r, fmt.Sprintf("/logs/%s/1", id))
		} else if page > pageCount {
			return razlink.RedirectView(r, fmt.Sprintf("/logs/%s/%d", id, pageCount))
		}
		first = (page - 1) * logsPerPage
		last = (page * logsPerPage) - 1
	}

	logs, _ = db.GetLogs(id, first, last)
	return view(newLogPageData(id, logs, pageCount), &title)
}

func handleLogClear(db *razlink.DB, r *http.Request) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView(r, "Not found", http.StatusNotFound)
	}

	cookie, _ := r.Cookie(id)
	if cookie == nil || cookie.Value != e.PasswordHash {
		return razlink.RedirectView(r, "/logs-auth/"+id)
	}

	db.DeleteLogs(id)
	return razlink.RedirectView(r, "/logs/"+id+"/1")
}

func handleLogChangePwPage(db *razlink.DB, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)
	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView(r, "Not found", http.StatusNotFound)
	}

	if r.Method != "POST" {
		return view(nil, nil)
	}

	r.ParseForm()
	oldpw := r.FormValue("old_password")
	pw := r.FormValue("password")
	pw2 := r.FormValue("confirm_password")

	if pw != pw2 {
		return view("Password mismatch", nil)
	}

	if !e.MatchPassword(oldpw) {
		return view("Wrong old password", nil)
	}

	e.SetPassword(pw)
	err := db.SetEntry(id, e)
	if err != nil {
		return razlink.ErrorView(r, err.Error(), http.StatusInternalServerError)
	}

	cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
	return razlink.CookieAndRedirectView(r, cookie, "/logs/"+id)
}

// GetLogPages ...
func GetLogPages(db *razlink.DB, logsPerPage int) []*razlink.Page {
	return []*razlink.Page{
		{
			Path:            "/logs/",
			Title:           "Logs",
			ContentTemplate: logPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogPage(db, logsPerPage, r, view)
			},
		},
		{
			Path:            "/logs-auth/",
			Title:           "Logs authentication",
			ContentTemplate: logAuthPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogAuthPage(db, r, view)
			},
		},
		{
			Path: "/logs-clear/",
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogClear(db, r)
			},
		},
		{
			Path:            "/logs-change-password/",
			Title:           "Change password",
			ContentTemplate: logChangePwPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleLogChangePwPage(db, r, view)
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
