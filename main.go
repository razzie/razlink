package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

var addPage = `
<form method="post">
	URL: <input type="text" name="url" /><br />
	Log password: <input type="password" name="password" /><br />
	<input type="radio" name="method" value="proxy" checked />Proxy
	<input type="radio" name="method" value="redirect" />Redirect<br />
	<input type="submit" value="Submit" />
</form>
`

var addResultPage = `
<a href="http://{{.Hostname}}/x/{{.ID}}">{{.Hostname}}/x/{{.ID}}</a><br />
<a href="http://{{.Hostname}}/logs/{{.ID}}">{{.Hostname}}/logs/{{.ID}}</a>
`

var logsPasswordPage = `
<form method="post">
	<input type="password" name="password" />
	<input type="submit" value="Submit" />
</form>
`

var logsPage = `
<ul>
	{{range .Logs}}
	<li>{{.}}</li>
	{{end}}
</ul>
`

func main() {
	hostname := flag.String("hostname", "link.gorzsony.com", "Hostname")
	port := flag.Int("port", 8081, "Port")
	flag.Parse()

	db := NewDB()

	addResultPageT, err := template.New("").Parse(addResultPage)
	if err != nil {
		panic(err)
	}

	logsPageT, err := template.New("").Parse(logsPage)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, addPage)
			return
		}

		r.ParseForm()
		url := r.FormValue("url")
		pw := r.FormValue("password")
		proxy := r.FormValue("method") == "proxy"

		e := db.InsertEntry(url, pw, proxy)
		http.Redirect(w, r, "/add/"+e.ID, http.StatusSeeOther)
	})

	http.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		id := r.URL.Path[5:]
		view := struct {
			Hostname string
			ID       string
		}{
			Hostname: *hostname,
			ID:       id,
		}
		addResultPageT.Execute(w, view)
	})

	http.HandleFunc("/x/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[3:]
		e := db.GetEntry(id)
		if e == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if e.Proxy {
			req, _ := http.NewRequest("GET", e.URL, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for k, v := range resp.Header {
				w.Header().Set(k, v[0])
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			resp.Body.Close()
		} else {
			http.Redirect(w, r, e.URL, http.StatusSeeOther)
		}

		db.InsertLog(id, r.Header.Get("X-REAL-IP"))
	})

	http.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Method != "POST" {
			fmt.Fprint(w, logsPasswordPage)
			return
		}

		r.ParseForm()
		pw := r.FormValue("password")

		id := r.URL.Path[6:]
		e := db.GetEntry(id)
		if e == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if !e.MatchPassword(pw) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logsPageT.Execute(w, e)
	})

	http.ListenAndServe("localhost:"+strconv.Itoa(*port), nil)
}
