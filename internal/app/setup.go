package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ashishxcode/commit-chronicle/internal/collect"
	"github.com/ashishxcode/commit-chronicle/internal/config"
	"github.com/ashishxcode/commit-chronicle/internal/tui"
)

const manualLabel = "📁 Enter a path manually…"

// firstRunSetup is invoked when no repos are configured and we have a TTY. It
// finds likely repo locations, lets the user pick one, optionally remembers it
// for next time, and returns the resolved repo paths to scan now.
func firstRunSetup() ([]string, error) {
	fmt.Fprintln(os.Stderr, "👋 Setup — let's find your git repositories.")
	fmt.Fprintln(os.Stderr, "   (runs automatically on first use; re-run any time with `commit-chronicle --setup`)")
	fmt.Fprintln(os.Stderr)

	cands := config.ScanCommonRoots()

	var options []string
	for _, c := range cands {
		options = append(options, fmt.Sprintf("%s  (%d repo%s)", c.Path, c.Count, plural(c.Count)))
	}
	options = append(options, manualLabel)

	idx, canceled, err := tui.Choose("📂 Where are your repositories?", options)
	if err != nil {
		return nil, err
	}
	if canceled {
		return nil, fmt.Errorf("setup canceled")
	}

	var root string
	if idx < len(cands) {
		root = cands[idx].Path
	} else {
		root, err = promptPath("Path to a folder containing your repos (e.g. ~/work): ")
		if err != nil {
			return nil, err
		}
		if root == "" {
			return nil, fmt.Errorf("setup canceled")
		}
		if n := config.CountRepos(root); n == 0 {
			fmt.Fprintf(os.Stderr, "⚠️  no git repos found under %s — continuing anyway.\n", root)
		}
	}

	// Resolve the repos under the chosen root now.
	repos, err := config.ResolveRepos(nil, []string{root})
	if err != nil {
		return nil, fmt.Errorf("scanning %s: %w", root, err)
	}

	// Offer to remember the choice.
	saveIdx, canceled, err := tui.Choose(
		fmt.Sprintf("💾 Remember %s for next time?", root),
		[]string{"Yes — save it", "No — just this once"})
	if err == nil && !canceled && saveIdx == 0 {
		if path, err := config.SaveRoot(root); err == nil {
			fmt.Fprintf(os.Stderr, "✅ saved to %s\n", path)
		} else {
			fmt.Fprintf(os.Stderr, "⚠️  could not save config: %v\n", err)
		}
	}

	reportGHStatus()
	fmt.Fprintln(os.Stderr)
	return repos, nil
}

// promptPath reads a single line of input, trimmed. Empty means cancel.
func promptPath(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && line == "" {
		return "", nil
	}
	return strings.TrimSpace(line), nil
}

// reportGHStatus prints whether PR/review discovery is available via gh.
func reportGHStatus() {
	switch {
	case !collect.HasGH():
		fmt.Fprintln(os.Stderr, "ℹ️  gh (GitHub CLI) not found — commits only. Install it + run `gh auth login` to include PRs & reviews.")
	case ghLogin() == "":
		fmt.Fprintln(os.Stderr, "ℹ️  gh is installed but not authenticated — run `gh auth login` to include PRs & reviews.")
	default:
		fmt.Fprintln(os.Stderr, "✅ gh authenticated — PRs & reviews will be included.")
	}
}
