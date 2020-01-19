package razlink

import (
	"html/template"
	"net/http"
)

var layout = `
{{define "layout"}}
<!DOCTYPE html>
<html>
	<head>
		<title>{{.Title}}</title>
		<link rel="icon" href="favicon.svg" type="image/svg+xml" />
		<style>
		div.outer {
			display: flex;
			align-items: center;
			justify-content: center;
		}
		div.inner {
			border: 1px solid black;
			padding: 1rem;
			display: inline-flex;
		}
		</style>
	</head>
	<body>
		<div class="outer">
			<div class="inner">
				<div>
				{{template "page" .Data}}
				</div>
			</div>
		</div>
	</body>
</html>
{{end}}
`

// Layout ...
type Layout struct {
	tmpl *template.Template
}

// NewLayout creates a new Layout
func NewLayout() *Layout {
	return &Layout{
		tmpl: template.Must(template.New("").Parse(layout)),
	}
}

// CreatePageRenderer creates a page renderer function
func (layout *Layout) CreatePageRenderer(title, content string, requestToData PageHandler) (func(http.ResponseWriter, *http.Request), error) {
	clone, _ := layout.tmpl.Clone()
	tmpl, err := clone.Parse(content)
	if err != nil {
		return nil, err
	}

	if requestToData == nil {
		requestToData = func(w http.ResponseWriter, r *http.Request) interface{} {
			return ""
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := requestToData(w, r)
		if data == nil {
			return
		}

		view := struct {
			Title string
			Data  interface{}
		}{
			Title: title,
			Data:  data,
		}

		tmpl.ExecuteTemplate(w, "layout", view)
	}, nil
}
