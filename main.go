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

	addr := "localhost:" + strconv.Itoa(*port)
	srv := newServer(db, *hostname)
	http.ListenAndServe(addr, srv)
}
