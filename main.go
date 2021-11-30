package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func parseBranches(raw []byte) []string {
	str := string(raw[:])
	split := strings.Split(str, "\n")

	var cleaned []string
	for _, v := range split {
		replaced := strings.ReplaceAll(v, "*", "")
		trimmed := strings.TrimSpace(replaced)
		if len(trimmed) > 0 {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

func hasRemote(remote string, branch string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", remote, branch)

	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, nil
	}
	str := string(result[:])
	return len(str) > 0, nil
}

func getRemote() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	result, err := cmd.CombinedOutput()
	parsed := string(result[:])
	return strings.TrimSpace(parsed), err
}

func deleteBranch(branch string) {
	cmd := exec.Command("git", "branch", "-D", branch)
	cmd.Run()
}

func main() {

	var (
		noop = func(context.Context, []string) error { return flag.ErrHelp }
	)

	globalFlags := flag.NewFlagSet("gitpurge", flag.ExitOnError)
	verboseFlag := globalFlags.Bool("verbose", false, "Verbose logging")

	list := &ffcli.Command{
		Name:       "list",
		ShortUsage: "gitpurge list",
		ShortHelp:  "list branches",
		LongHelp:   "List all branches and details about the branch",
		Exec: func(ctx context.Context, args []string) error {

			remote, err := getRemote()
			if err != nil {
				return err
			}

			cmd := exec.Command("git", "branch")
			result, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}
			branches := parseBranches(result)

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Name", "Remote Exists"})

			var rows []table.Row
			for _, branch := range branches {

				remoteExists, err := hasRemote(remote, branch)
				if err != nil {
					rows = append(rows, table.Row{branch, "unknown"})
				} else {
					strStatus := "N"
					if remoteExists {
						strStatus = "Y"
					}
					rows = append(rows, table.Row{branch, strStatus})
				}

			}

			t.AppendRows(rows)
			t.Render()

			return nil
		},
	}

	purge := &ffcli.Command{
		Name:       "purge",
		ShortUsage: "gitpurge purge",
		ShortHelp:  "delete local branches with no remote",
		LongHelp:   "Will delete all branches locally that do have a remote branch",
		Exec: func(ctx context.Context, args []string) error {

			remote, err := getRemote()
			if err != nil {
				return err
			}

			cmd := exec.Command("git", "branch")
			result, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}

			branches := parseBranches(result)
			for _, branch := range branches {

				remoteExists, err := hasRemote(remote, branch)
				if err == nil && !remoteExists {
					if *verboseFlag {
						fmt.Printf("Parsing Branch: %v, remoteExists: %v, deleting \n", branch, remoteExists)
					}
					deleteBranch(branch)
				} else {
					if *verboseFlag {
						fmt.Printf("Parsing Branch: %v, remoteExists: %v, not deleting \n", branch, remoteExists)
					}
				}
			}

			return nil
		},
	}

	gitpurge := &ffcli.Command{
		ShortUsage:  "gitpurge [global flags] <subcommand> [subcommand flags] [subcommand args]",
		ShortHelp:   "tool to help cleanup local branches",
		Subcommands: []*ffcli.Command{list, purge},
		FlagSet:     globalFlags,
		Exec:        noop,
	}

	if err := gitpurge.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "mapi: %v\n", err)
		}
		os.Exit(1)
	}
}
