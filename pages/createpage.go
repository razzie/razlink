package pages

import (
	"net/http"
	"net/url"
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
{{if .ShorthandURL}}
	Original URL:<br />
	<a href="{{.URL}}">{{.ShorthandURL}}</a><br />
	<br />
{{end}}
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

func handleCreatePage(db *razlink.DB, pr *razlink.PageRequest) *razlink.View {
	r := pr.Request
	if r.Method != "POST" {
		return nil
	}

	r.ParseForm()
	url := r.FormValue("url")
	pw := r.FormValue("password")
	pw2 := r.FormValue("confirm_password")

	if pw != pw2 {
		return pr.Respond("Password mismatch")
	}

	if url != "." && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	method, err := razlink.GetServeMethodForURL(r.Context(), url, time.Second*3)
	if err != nil {
		return pr.Respond(err.Error(), razlink.WithError(err, http.StatusInternalServerError))
	}

	e := razlink.NewEntry(url, pw, method)
	id, err := db.InsertEntry(nil, e)
	if err != nil {
		return pr.ErrorView(err.Error(), http.StatusInternalServerError)
	}

	db.InsertLog(id, r)

	cookie := &http.Cookie{Name: id, Value: e.PasswordHash, Path: "/"}
	return razlink.CookieAndRedirectView(r, cookie, "/link/"+id)
}

func handleCreateResultPage(db *razlink.DB, pr *razlink.PageRequest) *razlink.View {
	id, _ := getIDFromRequest(pr)

	e, _ := db.GetEntry(id)
	if e == nil {
		return pr.ErrorView("Not found", http.StatusNotFound)
	}

	data := struct {
		URL          string
		ShorthandURL string
		Hostname     string
		ID           string
		Decoy        string
		Track        bool
	}{
		Hostname: pr.Request.Host,
		ID:       id,
		Track:    e.Method == razlink.Track,
	}

	cookie, _ := pr.Request.Cookie(id)
	if cookie != nil && cookie.Value == e.PasswordHash {
		data.URL = e.URL

		if e.URL != "." {
			u, _ := url.Parse(e.URL)
			if u != nil {
				data.ShorthandURL = u.Host + razlink.GetShorthandPath(u.Path)
			}
		}

		filename := filepath.Base(e.URL)
		if len(filename) < 2 {
			data.Decoy = ""
		} else {
			data.Decoy = razlink.GetShorthandPath(filename)
		}
	}

	pr.Title = "Link: " + id
	return pr.Respond(&data)
}

// GetCreatePages ...
func GetCreatePages(db *razlink.DB) []*razlink.Page {
	return []*razlink.Page{
		{
			Path:            "/create",
			Title:           "Create a new link",
			ContentTemplate: createPageT,
			Handler: func(pr *razlink.PageRequest) *razlink.View {
				return handleCreatePage(db, pr)
			},
		},
		{
			Path:            "/link/",
			ContentTemplate: createResultPageT,
			Handler: func(pr *razlink.PageRequest) *razlink.View {
				return handleCreateResultPage(db, pr)
			},
		},
		{
			Path: "/add/", // for legacy bookmarks
			Handler: func(pr *razlink.PageRequest) *razlink.View {
				id, _ := getIDFromRequest(pr)
				return pr.RedirectView("/link/" + id)
			},
		},
	}
}
