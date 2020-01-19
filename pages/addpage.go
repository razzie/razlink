package pages

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/razzie/razlink"
)

var addPageT = `
{{define "page"}}
<form method="post">
	URL:<br />
	<input type="text" name="url" /><br />
	<br />
	Log password:<br />
	<input type="password" name="password" /><br />
	<br />
	<input type="submit" value="Submit" />
</form>
{{end}}
`

var addResultPageT = `
{{define "page"}}
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
{{end}}
`

func handleAddPage(db *razlink.DB, w http.ResponseWriter, r *http.Request) interface{} {
	if r.Method != "POST" {
		return ""
	}

	r.ParseForm()
	url := r.FormValue("url")
	pw := r.FormValue("password")
	method, err := razlink.GetServeMethodForURL(r.Context(), url, time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	e := razlink.NewEntry(url, pw, method)
	id, err := db.InsertEntry(nil, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	db.InsertLog(id, r)

	http.SetCookie(w, &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"})
	http.Redirect(w, r, "/add/"+id, http.StatusSeeOther)
	return nil
}

func handleAddResultPage(db *razlink.DB, hostname string, w http.ResponseWriter, r *http.Request) interface{} {
	id, _ := getIDFromRequest(r)

	e, _ := db.GetEntry(id)
	if e == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return nil
	}

	decoy := filepath.Base(e.URL)
	if len(decoy) < 2 {
		decoy = ""
	}

	view := struct {
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

	return &view
}

// GetAddPages ...
func GetAddPages(db *razlink.DB, hostname string) []*razlink.Page {
	return []*razlink.Page{
		&razlink.Page{
			Path:            "/add",
			Title:           "Create a new link",
			ContentTemplate: addPageT,
			Handler: func(w http.ResponseWriter, r *http.Request) interface{} {
				return handleAddPage(db, w, r)
			},
		},
		&razlink.Page{
			Path:            "/add/",
			Title:           "Bookmark this page!",
			ContentTemplate: addResultPageT,
			Handler: func(w http.ResponseWriter, r *http.Request) interface{} {
				return handleAddResultPage(db, hostname, w, r)
			},
		},
	}
}
