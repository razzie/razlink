package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/razzie/razlink"
)

// CLI ...
type CLI struct {
	cmds map[string]func(args []string)
}

// NewCLI ...
func NewCLI(db *razlink.DB) *CLI {
	cli := &CLI{
		cmds: make(map[string]func(args []string)),
	}

	cli.cmds["help"] = func(args []string) {
		for cmd := range cli.cmds {
			fmt.Println(cmd)
		}
	}

	cli.cmds["links"] = func(args []string) {
		if len(args) != 1 {
			fmt.Println("usage: links <pattern>")
			return
		}

		pattern := args[0]
		entries, err := db.GetEntries(pattern)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		for _, e := range entries {
			fmt.Println(e.ID, "-", e.URL)
		}
	}

	cli.cmds["add"] = func(args []string) {
		if len(args) != 2 {
			fmt.Println("usage: add <URL> <password>")
			return
		}

		url := args[0]
		pw := args[1]
		method, err := razlink.GetServeMethodForURL(context.Background(), url)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e, err := db.InsertEntry(url, pw, method)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("added:", e.ID)
	}

	cli.cmds["add-permanent"] = func(args []string) {
		if len(args) != 3 {
			fmt.Println("usage: add-permament <ID> <URL> <password>")
			return
		}

		ID := args[0]
		url := args[1]
		pw := args[2]
		method, err := razlink.GetServeMethodForURL(context.Background(), url)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e, err := db.InsertPermanentEntry(ID, url, pw, method)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("added:", e.ID)
	}

	cli.cmds["delete"] = func(args []string) {
		if len(args) != 1 {
			fmt.Println("usage: delete <ID>")
			return
		}

		ID := args[0]
		err := db.DeleteEntry(ID)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("deleted:", ID)
	}

	cli.cmds["logs"] = func(args []string) {
		if len(args) != 3 {
			fmt.Println("usage: logs <ID> <first> <last>")
			return
		}

		ID := args[0]
		first, _ := strconv.Atoi(args[1])
		last, _ := strconv.Atoi(args[2])
		logs, err := db.GetLogs(ID, first, last)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		for _, log := range logs {
			fmt.Println(log.String())
		}
	}

	cli.cmds["clear-logs"] = func(args []string) {
		if len(args) != 1 {
			fmt.Println("usage: clear-logs <ID>")
			return
		}

		ID := args[0]
		err := db.DeleteLogs(ID)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("logs cleared:", ID)
	}

	cli.cmds["exit"] = func(args []string) {
		os.Exit(0)
	}

	return cli
}

func (cli *CLI) handleCommand(cmd string) {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return
	}

	fn, ok := cli.cmds[args[0]]
	if !ok {
		fmt.Println("unknown command:", cmd)
		return
	}

	fn(args[1:])
}

// Run ...
func (cli *CLI) Run() {
	fmt.Println("awesome razlink command line")
	fmt.Println("----------------------------")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSuffix(cmd, "\n")

		cli.handleCommand(cmd)
	}
}
