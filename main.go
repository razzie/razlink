package main

import (
	"flag"
	"net/http"
	"strconv"
)

func main() {
	hostname := flag.String("hostname", "localhost:8080", "Hostname")
	port := flag.Int("port", 8080, "Port")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis hostname:port")
	redisPw := flag.String("redis-pw", "", "Redis password")
	redisDb := flag.Int("redis-db", 0, "Redis database (0-15)")
	flag.Parse()

	db, err := NewDB(*redisAddr, *redisPw, *redisDb)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	mux := http.DefaultServeMux
	installAddPage(db, mux, *hostname)
	installViewPage(db, mux)
	installLogPage(db, mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/add", http.StatusSeeOther)
	})

	http.ListenAndServe("localhost:"+strconv.Itoa(*port), mux)
}
