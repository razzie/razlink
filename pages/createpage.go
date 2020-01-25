package pages

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/razzie/razlink"
)

var createPageT = `
{{if .}}
<strong style="color: red">{{.}}</strong><br /><br />
{{end}}
<form method="post">
	<input type="text" name="url" placeholder="URL" style="min-width: 400px" /><br />
	<input type="password" name="password" placeholder="Password for logs" /><br />
	<input type="password" name="confirm_password" placeholder="Confirm password" /><br />
	<button>Create</button>
</form>
`

var createResultPageT = `
<strong>Bookmark this page!</strong><br />
<br />
{{if .Track}}
	Embed this in your website:<br />
	&lt;img src="<a href="http://{{.Hostname}}/x/{{.ID}}">http://{{.Hostname}}/x/{{.ID}}</a>" width="1" height="1" /&gt;<br />
{{else}}
	Access the target URL:<br />
	<a href="http://{{.Hostname}}/x/{{.ID}}">{{.Hostname}}/x/{{.ID}}</a><br />
	{{if .Decoy}}
		<a href="http://{{.Hostname}}/x/{{.ID}}/{{.Decoy}}">{{.Hostname}}/x/{{.ID}}/{{.Decoy}}</a><br />
	{{end}}
{{end}}
<br />
Access logs:<br />
<a href="http://{{.Hostname}}/logs/{{.ID}}">{{.Hostname}}/logs/{{.ID}}</a>
`

func handleCreatePage(db *razlink.DB, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	if r.Method != "POST" {
		return view(nil, nil)
	}

	r.ParseForm()
	url := r.FormValue("url")
	pw := r.FormValue("password")
	pw2 := r.FormValue("confirm_password")

	if pw != pw2 {
		return view("Password mismatch", nil)
	}

	if url != "." && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	method, err := razlink.GetServeMethodForURL(r.Context(), url, time.Second)
	if err != nil {
		return razlink.ErrorView(r, err.Error(), http.StatusInternalServerError)
	}

	e := razlink.NewEntry(url, pw, method)
	id, err := db.InsertEntry(nil, e)
	if err != nil {
		return razlink.ErrorView(r, err.Error(), http.StatusInternalServerError)
	}

	db.InsertLog(id, r)

	cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
	return razlink.CookieAndRedirectView(r, cookie, "/link/"+id)
}

func handleCreateResultPage(db *razlink.DB, hostname string, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)

	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView(r, "Not found", http.StatusNotFound)
	}

	decoy := filepath.Base(e.URL)
	if len(decoy) < 2 {
		decoy = ""
	}

	data := struct {
		Hostname string
		ID       string
		Decoy    string
		Track    bool
	}{
		Hostname: hostname,
		ID:       id,
		Decoy:    decoy,
		Track:    e.Method == razlink.Track,
	}

	return view(&data, &id)
}

// GetCreatePages ...
func GetCreatePages(db *razlink.DB, hostname string) []*razlink.Page {
	return []*razlink.Page{
		&razlink.Page{
			Path:            "/create",
			Title:           "Create a new link",
			ContentTemplate: createPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleCreatePage(db, r, view)
			},
		},
		&razlink.Page{
			Path:            "/link/",
			ContentTemplate: createResultPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleCreateResultPage(db, hostname, r, view)
			},
		},
		&razlink.Page{
			Path: "/add/", // for legacy bookmarks
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				id, _ := getIDFromRequest(r)
				return razlink.RedirectView(r, "/link/"+id)
			},
		},
	}
}