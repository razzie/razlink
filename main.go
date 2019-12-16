package main

import (
	"flag"
	"net/http"
	"strconv"
)

func main() {
	hostname := flag.String("hostname", "link.gorzsony.com", "Hostname")
	port := flag.Int("port", 8081, "Port")
	flag.Parse()

	db := NewDB("localhost:6379", "", 0)

	mux := http.DefaultServeMux
	installAddPage(db, mux, *hostname)
	installViewPage(db, mux)
	installLogPage(db, mux)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/add", http.StatusSeeOther)
	})

	http.ListenAndServe("localhost:"+strconv.Itoa(*port), nil)
}
