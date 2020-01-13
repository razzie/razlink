package main

import (
	"fmt"
	"net/http"
)

var welcomePage = `
<div style="display: flex; align-items: center; justify-content: center">
	<div style="border: 1px solid black; padding: 1rem; display: inline-flex">
		<div style="display: block">
			<strong>Welcome to razlink!</strong><br />
			<br />
			This is a lightweight link proxy/redirect service with logging.<br />
			<br />
			You can <a href="add">add a new link</a> with a custom password to view the visits.<br />
			There are 3 automatic operation modes based on the URL:<br />
			<ul>
				<li>Proxy - for files (e.g. an image)</li>
				<li>Embed - for websites that allow embedding</li>
				<li>Redirect - for websites that forbid embedding</li>
			</ul><br />
			Check out the source code at <a href="https://github.com/razzie/razlink" target="_blank">github.com/razzie/razlink</a>.<br />
		</div>
	</div>
</div>
`

func installWelcomePage(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, welcomePage)
	})
}
