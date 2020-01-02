package main

import (
	"fmt"
	"net/http"
)

type server struct {
	mux http.ServeMux
}

func newServer(db *DB, hostname string) *server {
	srv := &server{}
	mux := &srv.mux

	installAddPage(db, mux, hostname)
	installViewPage(db, mux)
	installLogPage(db, mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/add", http.StatusSeeOther)
	})

	return srv
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go fmt.Println(NewLog(r.Header.Get("X-REAL-IP"), r.UserAgent()), "- path:", r.URL.Path)
	srv.mux.ServeHTTP(w, r)
}
