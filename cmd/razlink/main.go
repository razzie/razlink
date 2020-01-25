package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/razzie/razlink"
	"github.com/razzie/razlink/pages"
)

func main() {
	hostname := flag.String("hostname", "localhost:8080", "Hostname")
	port := flag.Int("port", 8080, "Port")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis hostname:port")
	redisPw := flag.String("redis-pw", "", "Redis password")
	redisDb := flag.Int("redis-db", 0, "Redis database (0-15)")
	cliMode := flag.Bool("cli", false, "Enable CLI mode instead of http server")
	viewMode := flag.Bool("view-mode", false, "View-mode disables welcome and create pages")
	flag.Parse()

	db, err := razlink.NewDB(*redisAddr, *redisPw, *redisDb)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if *cliMode {
		NewCLI(db).Run()
	} else {
		fmt.Println("Starting Razlink instance:", razlink.InstanceID)
		addr := "localhost:" + strconv.Itoa(*port)
		srv := razlink.NewServer()
		if !*viewMode {
			srv.AddPages(append(pages.GetCreatePages(db, *hostname), pages.GetWelcomePage())...)
		}
		srv.AddPages(append(pages.GetLogPages(db, 20), pages.GetViewPage(db))...)
		http.ListenAndServe(addr, srv)
	}
}
