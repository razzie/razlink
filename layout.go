package razlink

import (
	"html/template"
	"net/http"
)

var styleT = `
a {
	color: black;
	text-decoration: underline;
	text-decoration-color: rgb(220, 53, 69);
	-webkit-text-decoration-color: rgb(220, 53, 69);
}
a:hover {
	color: dimgrey;
}
input[type="text"], input[type="password"] {
	border: 0;
	outline: 0;
	background: transparent;
	border-bottom: 1px solid black;
	margin-bottom: 1rem;
	min-width: 250px;
}
input[type="submit"], input[type="button"], button, .button {
	border: 1px solid black;
	border-radius: 5px;
	background-color: whitesmoke;
	padding: 5px 10px;
	margin: 10px 0;
	color: black;
	text-decoration: none;
	display: inline-block;
	cursor: pointer;
}
input[type="submit"]:disabled, input[type="button"]:disabled, button:disabled, .button:disabled {
	color: grey;
}
input[type="submit"]:not(:disabled):hover, input[type="button"]:not(:disabled):hover,
button:not(:disabled):hover, .button:not(:disabled):hover {
	background-color: lightsteelblue;
}
table {
	border-collapse: collapse;
	margin-bottom: 1rem;
	border-spacing: 0;
}
td {
	padding: 10px;
	border: 1px solid transparent;
}
tr:nth-child(odd) > td {
	background-color: whitesmoke;
}
tr:first-child > td {
	font-weight: bold;
	border-bottom: 1px solid black;
	background-color: white;
}
tr:not(:first-child):hover > td {
	background-color: lightsteelblue;
}
tr:not(:first-child) > td:first-child {
	border-radius: 10px 0 0 10px;
	border: 0;
}
tr:not(:first-child) > td:last-child {
	border-radius: 0 10px 10px 0;
	border: 0;
}
small {
	color: dimgrey;
}
`

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
				border-radius: 15px;
			}
			{{template "style"}}
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

// Layout is used to give pages a uniform layout
type Layout interface {
	BindTemplate(pageTemplate string, stylesheets, scripts []string, meta map[string]string) (LayoutRenderer, error)
}

// LayoutRenderer is a function that renders a html page
type LayoutRenderer func(w http.ResponseWriter, r *http.Request, title string, data interface{}, statusCode int)

// DefaultLayout is razlink's default layout
var DefaultLayout Layout = (*layout)(template.Must(template.New("layout").Parse(layoutT)))

type layout template.Template

// BindTemplate creates a layout renderer function from a page template
func (l *layout) BindTemplate(pageTemplate string, stylesheets, scripts []string, meta map[string]string) (LayoutRenderer, error) {
	cloneLayout, _ := (*template.Template)(l).Clone()
	tmpl, err := cloneLayout.New("page").Parse(pageTemplate)
	if err != nil {
		return nil, err
	}

	if tmpl.Lookup("style") == nil {
		tmpl = template.Must(tmpl.New("style").Parse(styleT))
	}

	return func(w http.ResponseWriter, r *http.Request, title string, data interface{}, statusCode int) {
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

		w.WriteHeader(statusCode)
		tmpl.ExecuteTemplate(w, "layout", &view)
	}, nil
}
