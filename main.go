package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {

	var (
		noop = func(context.Context, []string) error { return flag.ErrHelp }
	)

	list := &ffcli.Command{
		Name:       "list",
		ShortUsage: "gitpurge list",
		ShortHelp:  "list branches with no remote",
		LongHelp:   "List all branches that have no remote branch",
		Exec: func(ctx context.Context, args []string) error {
			cmd := exec.Command("git", "branch")
			result, err := cmd.CombinedOutput()
			fmt.Println(string(result[:]))
			fmt.Println(err)
			return nil
		},
	}

	purge := &ffcli.Command{
		ShortUsage:  "gitpurge [global flags] <subcommand> [subcommand flags] [subcommand args]",
		ShortHelp:   "tool to help cleanup local branches",
		Subcommands: []*ffcli.Command{list},
		Exec:        noop,
	}

	if err := purge.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "mapi: %v\n", err)
		}
		os.Exit(1)
	}
}
