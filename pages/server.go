package pages

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/razzie/razlink"
)

// Server ...
type Server struct {
	mux http.ServeMux
}

// NewServer ...
func NewServer(db *razlink.DB, hostname string) *Server {
	srv := &Server{}
	mux := &srv.mux

	installAddPage(db, mux, hostname)
	installViewPage(db, mux)
	installLogPage(db, mux)
	installWelcomePage(mux)

	/*mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/add", http.StatusSeeOther)
	})*/

	return srv
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
	fmt.Println(razlink.NewLog(r))
}

func getIDFromRequest(r *http.Request) (string, bool) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return "", false
	}
	return parts[2], true
}
