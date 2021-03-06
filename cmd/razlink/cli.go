package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
		if len(args) > 1 {
			fmt.Println("usage: links <pattern>")
			return
		} else if len(args) == 0 {
			args = []string{"*"}
		}

		pattern := args[0]
		entries, err := db.GetEntries(pattern)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No results")
			return
		}

		fmt.Println("Permanent entries are marked with *")

		for id, e := range entries {
			if e.Permanent {
				id += "*"
			}
			fmt.Println(id, "-", e.URL)
		}
	}

	cli.cmds["add"] = func(args []string) {
		if len(args) != 2 {
			fmt.Println("usage: add <URL> <password>")
			return
		}

		url := args[0]
		pw := args[1]
		method, err := razlink.GetServeMethodForURL(context.Background(), url, time.Minute)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e := razlink.NewEntry(url, pw, method)
		id, err := db.InsertEntry(nil, e)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("added:", id)
	}

	cli.cmds["add-permanent"] = func(args []string) {
		if len(args) != 3 {
			fmt.Println("usage: add-permament <ID> <URL> <password>")
			return
		}

		id := args[0]
		url := args[1]
		pw := args[2]
		method, err := razlink.GetServeMethodForURL(context.Background(), url, time.Minute)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e := razlink.NewEntry(url, pw, method)
		e.Permanent = true
		_, err = db.InsertEntry(&id, e)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("added:", id)
	}

	cli.cmds["change-password"] = func(args []string) {
		if len(args) != 2 {
			fmt.Println("usage: change-password <ID> <new password>")
			return
		}

		id := args[0]
		pw := args[1]
		e, err := db.GetEntry(id)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e.SetPassword(pw)
		err = db.SetEntry(id, e)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("done")
	}

	cli.cmds["set-permanent"] = func(args []string) {
		if len(args) != 2 {
			fmt.Println("usage: set-permanent <ID> <true/false>")
			return
		}

		id := args[0]
		perm, err := strconv.ParseBool(args[1])
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e, err := db.GetEntry(id)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		e.Permanent = perm
		err = db.SetEntry(id, e)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("done")
	}

	cli.cmds["delete"] = func(args []string) {
		if len(args) != 1 {
			fmt.Println("usage: delete <ID>")
			return
		}

		id := args[0]
		err := db.DeleteEntry(id)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("deleted:", id)
	}

	cli.cmds["logs"] = func(args []string) {
		if len(args) == 1 {
			args = append(args, "0", "19")
		} else if len(args) != 3 {
			fmt.Println("usage: logs <ID> <first> <last>")
			return
		}

		id := args[0]
		first, _ := strconv.Atoi(args[1])
		last, _ := strconv.Atoi(args[2])
		logs, err := db.GetLogs(id, first, last)
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

		id := args[0]
		err := db.DeleteLogs(id)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("logs cleared:", id)
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
