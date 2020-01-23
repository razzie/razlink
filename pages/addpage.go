package pages

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/razzie/razlink"
)

var addPageT = `
<form method="post">
	URL:<br />
	<input type="text" name="url" /><br />
	<br />
	Log password:<br />
	<input type="password" name="password" /><br />
	<br />
	<input type="submit" value="Submit" />
</form>
`

var addResultPageT = `
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

func handleAddPage(db *razlink.DB, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	if r.Method != "POST" {
		return view(nil)
	}

	r.ParseForm()
	url := r.FormValue("url")
	pw := r.FormValue("password")

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	method, err := razlink.GetServeMethodForURL(r.Context(), url, time.Second)
	if err != nil {
		return razlink.ErrorView(err.Error(), http.StatusInternalServerError)
	}

	e := razlink.NewEntry(url, pw, method)
	id, err := db.InsertEntry(nil, e)
	if err != nil {
		return razlink.ErrorView(err.Error(), http.StatusInternalServerError)
	}

	db.InsertLog(id, r)

	cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
	return razlink.CookieAndRedirectView(r, cookie, "/add/"+id)
}

func handleAddResultPage(db *razlink.DB, hostname string, r *http.Request, view razlink.ViewFunc) razlink.PageView {
	id, _ := getIDFromRequest(r)

	e, _ := db.GetEntry(id)
	if e == nil {
		return razlink.ErrorView("Not found", http.StatusNotFound)
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

	return view(&data)
}

// GetAddPages ...
func GetAddPages(db *razlink.DB, hostname string) []*razlink.Page {
	return []*razlink.Page{
		&razlink.Page{
			Path:            "/add",
			Title:           "Create a new link",
			ContentTemplate: addPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleAddPage(db, r, view)
			},
		},
		&razlink.Page{
			Path:            "/add/",
			Title:           "Bookmark this page!",
			ContentTemplate: addResultPageT,
			Handler: func(r *http.Request, view razlink.ViewFunc) razlink.PageView {
				return handleAddResultPage(db, hostname, r, view)
			},
		},
	}
}
