package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/razzie/razlink"
	"github.com/razzie/razlink/pages"
)

// Command-line args
var (
	Port      int
	RedisAddr string
	RedisPw   string
	RedisDB   int
	CLIMode   bool
	ViewMode  bool
)

func init() {
	flag.IntVar(&Port, "port", 8080, "Port")
	flag.StringVar(&RedisAddr, "redis-addr", "localhost:6379", "Redis hostname:port")
	flag.StringVar(&RedisPw, "redis-pw", "", "Redis password")
	flag.IntVar(&RedisDB, "redis-db", 0, "Redis database (0-15)")
	flag.BoolVar(&CLIMode, "cli", false, "Enable CLI mode instead of http server")
	flag.BoolVar(&ViewMode, "view-mode", false, "View-mode disables welcome and create pages")
}

func main() {
	flag.Parse()

	db, err := razlink.NewDB(RedisAddr, RedisPw, RedisDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if CLIMode {
		NewCLI(db).Run()
	} else {
		fmt.Println("Starting Razlink instance:", razlink.InstanceID)
		addr := ":" + strconv.Itoa(Port)
		srv := razlink.NewServer()
		if !ViewMode {
			srv.AddPages(append(pages.GetCreatePages(db), pages.GetWelcomePage())...)
		}
		srv.AddPages(append(pages.GetLogPages(db, 20), pages.GetViewPage(db))...)
		log.Fatal(http.ListenAndServe(addr, srv))
	}
}
