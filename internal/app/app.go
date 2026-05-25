// Package app wires together config, collection, rendering and the TUI.
package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ashishxcode/commit-chronicle/internal/collect"
	"github.com/ashishxcode/commit-chronicle/internal/config"
	"github.com/ashishxcode/commit-chronicle/internal/model"
	"github.com/ashishxcode/commit-chronicle/internal/render"
	"github.com/ashishxcode/commit-chronicle/internal/tui"
	"github.com/mattn/go-isatty"
)

// Config is the fully-parsed CLI configuration.
type Config struct {
	Since, From, To, Month, Date string
	Author, User, Repos, Root    string
	Out, Format                  string
	NoEdit, All, Copy, NoPR      bool
}

// Run executes the whole pipeline.
func Run(c Config) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is required but was not found on PATH")
	}

	repos, err := config.ResolveRepos(splitCSV(c.Repos), splitCSV(c.Root))
	if err != nil {
		return err
	}

	author := c.Author
	if author == "" {
		author = gitConfigName(repos[0])
	}
	if author == "" {
		return fmt.Errorf("could not determine author; pass --author \"Your Name\"")
	}

	interactive := isTerminal()

	rng, err := resolveRange(c, interactive)
	if err != nil {
		return err
	}

	// GitHub login for PR/review discovery.
	ghUser := c.User
	if ghUser == "" && !c.NoPR && collect.HasGH() {
		ghUser = ghLogin()
	}

	opts := collect.Options{
		Repos:          repos,
		Author:         author,
		User:           ghUser,
		Range:          rng,
		IncludePRs:     !c.NoPR,
		IncludeReviews: !c.NoPR,
	}
	label := fmt.Sprintf(" scanning %d repo(s) for \"%s\" (%s) ", len(repos), author, rng.Label)

	var items []model.Item
	if interactive {
		// Animated spinner with live progress (Claude-style).
		items, err = tui.RunWithSpinner(label, func(report func(string, int)) ([]model.Item, error) {
			return collect.Gather(opts, report)
		})
	} else {
		fmt.Fprintf(os.Stderr, "🔎%s…\n", label)
		items, err = collect.Gather(opts, func(stage string, n int) {
			if n > 0 {
				fmt.Fprintf(os.Stderr, "   • %s: %d\n", stage, n)
			}
		})
	}
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("nothing found for \"%s\" in range (%s)", author, rng.Label)
	}

	// Pick
	selected := items
	if !c.All {
		if !interactive {
			return fmt.Errorf("no TTY for the picker; re-run with --all or in a terminal")
		}
		sel, canceled, err := tui.Pick(items, rng.Label, author)
		if err != nil {
			return err
		}
		if canceled || len(sel) == 0 {
			fmt.Fprintln(os.Stderr, "nothing selected — bye.")
			return nil
		}
		selected = sel
	}

	meta := render.Meta{Author: author, RangeLabel: rng.Label}

	var content string
	if c.Format == "json" {
		content = render.JSON(selected, meta)
	} else {
		content = render.Markdown(selected, meta)
		if !c.NoEdit && interactive {
			edited, canceled, err := tui.Edit(content)
			if err != nil {
				return err
			}
			if !canceled {
				content = edited
			}
		}
	}

	out := c.Out
	if out == "" {
		out = defaultOutPath(c.Format)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	if err := os.WriteFile(out, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", out, err)
	}
	fmt.Fprintf(os.Stderr, "✅ %d entr%s → %s\n", len(selected), plural(len(selected)), out)

	if c.Copy {
		if err := copyClipboard(content); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  clipboard: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "📋 copied to clipboard")
		}
	}

	if !interactive {
		fmt.Print(content)
	}
	return nil
}

func resolveRange(c Config, interactive bool) (model.Range, error) {
	switch {
	case c.Date != "":
		return model.FromDates(c.Date, c.Date)
	case c.Month != "":
		return model.FromMonth(c.Month)
	case c.From != "":
		return model.FromDates(c.From, c.To)
	case c.Since != "":
		return model.Range{Since: c.Since, Label: "since " + c.Since}, nil
	}
	if !interactive {
		return model.Range{Since: "30 days ago", Label: "last 30 days"}, nil
	}
	idx, canceled, err := tui.Choose("📅 Worklog period", model.PresetNames)
	if err != nil {
		return model.Range{}, err
	}
	if canceled {
		return model.Range{}, fmt.Errorf("canceled")
	}
	if model.PresetNames[idx] == model.CustomRangeLabel {
		from, to, canceled, err := tui.CustomRange()
		if err != nil {
			return model.Range{}, err
		}
		if canceled || from == "" {
			return model.Range{}, fmt.Errorf("canceled")
		}
		return model.FromDates(from, to)
	}
	return model.Preset(model.PresetNames[idx]), nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func gitConfigName(repo string) string {
	out, err := exec.Command("git", "-C", repo, "config", "user.name").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func ghLogin() string {
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func defaultOutPath(format string) string {
	ext := "md"
	if format == "json" {
		ext = "json"
	}
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "Downloads")
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		dir, _ = os.Getwd()
	}
	return filepath.Join(dir, fmt.Sprintf("commit-chronicle_%s.%s", time.Now().Format("20060102_150405"), ext))
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())
}

func copyClipboard(s string) error {
	var name string
	var args []string
	switch {
	case hasCmd("pbcopy"):
		name = "pbcopy"
	case hasCmd("wl-copy"):
		name = "wl-copy"
	case hasCmd("xclip"):
		name, args = "xclip", []string{"-selection", "clipboard"}
	default:
		return fmt.Errorf("no clipboard tool found (pbcopy/wl-copy/xclip)")
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(s)
	return cmd.Run()
}

func hasCmd(n string) bool { _, err := exec.LookPath(n); return err == nil }

func plural(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
