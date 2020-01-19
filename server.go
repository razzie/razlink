package razlink

import (
	"fmt"
	"net/http"
)

// Server ...
type Server struct {
	mux    http.ServeMux
	layout *Layout
}

// NewServer creates a new Server
func NewServer() *Server {
	srv := &Server{
		layout: NewLayout(),
	}

	srv.mux.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		WriteFavicon(w)
	})

	return srv
}

// AddPage adds a new servable page to the server
func (srv *Server) AddPage(page *Page) error {
	renderer, err := srv.layout.CreatePageRenderer(page.Title, page.ContentTemplate, page.Handler)
	if err != nil {
		return err
	}

	srv.mux.HandleFunc(page.Path, renderer)
	return nil
}

// AddPages adds multiple pages to the server and panics if anything goes wrong
func (srv *Server) AddPages(pages ...*Page) {
	for _, page := range pages {
		err := srv.AddPage(page)
		if err != nil {
			panic(err)
		}
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
	fmt.Println(NewLog(r), "- path:", r.URL.Path)
}
