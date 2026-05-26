// Command commit-chronicle: find your commits and PRs across repos in a date
// window, pick them interactively, and turn them into a worklog (markdown or
// json). A single self-contained binary — `git` is required, and `gh` is used
// opportunistically to discover PR and review activity.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ashishxcode/commit-chronicle/internal/app"
)

// version is set at build time via -ldflags "-X main.version=…".
var version = "dev"

func main() {
	for _, a := range os.Args[1:] {
		if a == "--version" || a == "-v" {
			fmt.Println("commit-chronicle " + version)
			return
		}
	}
	c, err := parseFlags()
	if err != nil {
		if err == flag.ErrHelp {
			return
		}
		fmt.Fprintln(os.Stderr, "❌ "+err.Error())
		os.Exit(2)
	}
	if err := app.Run(*c); err != nil {
		fmt.Fprintln(os.Stderr, "❌ "+err.Error())
		os.Exit(1)
	}
}

func parseFlags() (*app.Config, error) {
	c := &app.Config{}
	fs := flag.NewFlagSet("commit-chronicle", flag.ContinueOnError)
	fs.StringVar(&c.Since, "since", "", `relative range, e.g. "7 days ago"`)
	fs.StringVar(&c.From, "from", "", "start date YYYY-MM-DD (inclusive)")
	fs.StringVar(&c.To, "to", "", "end date YYYY-MM-DD (inclusive; default today)")
	fs.StringVar(&c.Month, "month", "", "whole calendar month YYYY-MM")
	fs.StringVar(&c.Date, "date", "", "single day YYYY-MM-DD")
	fs.StringVar(&c.Author, "author", "", "author to match (default: git config user.name)")
	fs.StringVar(&c.User, "user", "", "GitHub login for PR discovery (default: gh user)")
	fs.StringVar(&c.Repos, "repos", "", "comma-separated repo paths (overrides config)")
	fs.StringVar(&c.Root, "root", "", "comma-separated dirs to scan for git repos, e.g. ~/work")
	fs.StringVar(&c.Out, "out", "", "output path (default: Downloads, timestamped)")
	fs.StringVar(&c.Format, "format", "md", "output format: md | json")
	fs.BoolVar(&c.NoEdit, "no-edit", false, "skip the editor step")
	fs.BoolVar(&c.All, "all", false, "select everything (skip the picker)")
	fs.BoolVar(&c.NoPR, "no-pr", false, "skip GitHub PR + review discovery")
	fs.BoolVar(&c.Copy, "copy", false, "copy the whole worklog to the clipboard (skips the picker)")
	fs.BoolVar(&c.Setup, "setup", false, "re-run the guided repo setup")
	fs.Usage = usage
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	switch c.Format {
	case "md", "markdown":
		c.Format = "md"
	case "json":
	default:
		return nil, fmt.Errorf("unknown --format %q (use md or json)", c.Format)
	}
	return c, nil
}

func usage() {
	fmt.Fprint(os.Stderr, `📓 commit-chronicle — pick your commits & PRs, build a worklog (single binary)

USAGE:
    commit-chronicle [OPTIONS]

Pick a time range (interactively if no date flag), then gather everything you
did in that window — commits (git history + commits on PRs you authored),
PRs you authored, and PRs you reviewed — into one fuzzy-filterable, tagged
picker. Multi-select, optionally edit, then export to markdown/json.

OPTIONS:
    --since <when>      relative range, e.g. "7 days ago"
    --from YYYY-MM-DD   start date (inclusive)
    --to   YYYY-MM-DD   end date (inclusive; default today)
    --month YYYY-MM     whole calendar month
    --date  YYYY-MM-DD  single day
    --author "Name"     author to match (default: git config user.name)
    --user <login>      GitHub login for PR discovery (default: gh user)
    --repos a,b,c       comma-separated repo paths (overrides config)
    --root  ~/work      comma-separated dirs to auto-discover git repos under
    --out <path>        output path (default: Downloads, timestamped)
    --format md|json    output format (default: md)
    --all               select everything (skip the picker)
    --no-edit           skip the editor step
    --no-pr             skip GitHub PR + review discovery (git commits only)
    --copy              copy the whole worklog to the clipboard (skips the picker)
    --setup             re-run the guided repo setup (runs automatically on first use)
    -h, --help          show this help

REPO CONFIG (unioned): --repos · --root · ./.commit-chronicle ·
    ~/.config/commit-chronicle/repos · ~/.config/commit-chronicle/roots ·
    (fallback) the current dir if it's a git repo

PICKER KEYS: ↑↓ move · space/tab select · a all · / filter · enter confirm · q cancel

EXAMPLES:
    commit-chronicle --since "7 days ago"
    commit-chronicle --month 2026-05 --copy
    commit-chronicle --root ~/work --since "30 days ago"
    commit-chronicle --date 2026-05-25 --all --format json --out ./today.json
`)
}
