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
		<base href="{{.Base}}" target="_blank">
		<link rel="icon" href="favicon.svg" type="image/svg+xml" />
		<style>
		body {
			background-color: #ffffff;
			background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='4' height='4' viewBox='0 0 4 4'%3E%3Cpath fill='%23808080' fill-opacity='0.5' d='M1 3h1v1H1V3zm2-2h1v1H3V1z'%3E%3C/path%3E%3C/svg%3E");
		}
		div.outer {
			display: flex;
			align-items: center;
			justify-content: center;
		}
		div.inner {
			background-color: #ffffff;
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
			Base  string
			Data  interface{}
		}{
			Title: title,
			Base:  GetBase(r),
			Data:  data,
		}

		tmpl.ExecuteTemplate(w, "layout", view)
	}, nil
}
