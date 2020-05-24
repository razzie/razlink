package razlink

import (
	"html/template"
	"net/http"
)

var layoutT = `
<!DOCTYPE html>
<html>
	<head>
		{{if .Title}}<title>{{.Title}}</title>{{end}}
		<base href="{{.Base}}" />
		{{range $name, $content := .Meta}}
			<meta name="{{$name}}" content="{{$content}}" />
		{{end}}
		<link rel="icon" href="favicon.png" type="image/png" />
		<style>
		body {
			background-color: white;
			background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='4' height='4' viewBox='0 0 4 4'%3E%3Cpath fill='%23808080' fill-opacity='0.5' d='M1 3h1v1H1V3zm2-2h1v1H3V1z'%3E%3C/path%3E%3C/svg%3E");
		}
		a {
			color: black;
			text-decoration: underline;
			text-decoration-color: rgb(220, 53, 69);
			-webkit-text-decoration-color: rgb(220, 53, 69);
		}
		a:hover {
			color: dimgrey;
		}
		div.outer {
			display: flex;
			align-items: center;
			justify-content: center;
		}
		div.inner {
			background-color: white;
			border: 1px solid black;
			padding: 1rem;
			display: inline-flex;
		}
		input {
			border: 0;
			outline: 0;
			background: transparent;
			border-bottom: 1px solid black;
			margin-bottom: 1rem;
		}
		table {
			border-collapse: collapse;
			margin-bottom: 1rem;
		}
		tr:nth-child(odd) > td {
			background-color: whitesmoke;
		}
		tr:nth-child(1) > td {
			font-weight: bold;
			border-bottom: 1px solid black;
			background-color: white;
		}
		td {
			padding: 10px;
		}
		</style>
		{{range .Stylesheets}}
			<link rel="stylesheet" href="{{.}}" />
		{{end}}
		{{range .Scripts}}
			<script src="{{.}}"></script>
		{{end}}
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
`

var layout = template.Must(template.New("layout").Parse(layoutT))

// LayoutRenderer is a function that renders a html page
type LayoutRenderer func(w http.ResponseWriter, r *http.Request, title string, data interface{})

// BindLayout creates a layout renderer function
func BindLayout(pageTemplate string, stylesheets, scripts []string, meta map[string]string) (LayoutRenderer, error) {
	cloneLayout, _ := layout.Clone()
	tmpl, err := cloneLayout.New("page").Parse(pageTemplate)
	if err != nil {
		return nil, err
	}

	if meta == nil {
		meta = map[string]string{
			"author": "Gábor Görzsöny",
		}
	}

	return func(w http.ResponseWriter, r *http.Request, title string, data interface{}) {
		view := struct {
			Title       string
			Base        string
			Stylesheets []string
			Scripts     []string
			Meta        map[string]string
			Data        interface{}
		}{
			Title:       title,
			Base:        GetBase(r),
			Stylesheets: stylesheets,
			Scripts:     scripts,
			Meta:        meta,
			Data:        data,
		}

		tmpl.ExecuteTemplate(w, "layout", &view)
	}, nil
}
