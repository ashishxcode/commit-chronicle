// Package config resolves the list of git repositories to scan.
package config

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// skipDirs are never descended into during root discovery.
var skipDirs = map[string]bool{
	"node_modules": true, "vendor": true, ".Trash": true,
	"Library": true, ".cache": true, "dist": true, "build": true,
}

// ResolveRepos returns validated git repository paths from explicit paths,
// root directories to scan, and config files.
//
// Sources (unioned, then de-duplicated):
//   - explicit repo paths (e.g. --repos)
//   - repos discovered under root dirs (e.g. --root, or the `roots` config)
//   - ./.commit-chronicle                       (repo paths, one per line)
//   - $XDG_CONFIG_HOME/commit-chronicle/repos   (repo paths)
//   - $XDG_CONFIG_HOME/commit-chronicle/roots   (root dirs to scan)
//
// If none of those yield anything, the current directory is used when it is a
// git repo.
func ResolveRepos(explicit, roots []string) ([]string, error) {
	var candidates []string
	candidates = append(candidates, explicit...)

	// Root dirs: from the caller plus the roots config file.
	allRoots := append([]string{}, roots...)
	if rs, ok := linesFromFile(xdgPath("roots")); ok {
		allRoots = append(allRoots, rs...)
	}
	for _, root := range allRoots {
		candidates = append(candidates, discoverRepos(expand(root))...)
	}

	// Explicit repo-path config files.
	if c, ok := linesFromFile(".commit-chronicle"); ok {
		candidates = append(candidates, c...)
	}
	if c, ok := linesFromFile(xdgPath("repos")); ok {
		candidates = append(candidates, c...)
	}

	// Fallback: the current directory.
	if len(candidates) == 0 {
		cwd, err := os.Getwd()
		if err == nil && isGitRepo(cwd) {
			candidates = []string{cwd}
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no repos found\n" +
			"  configure via --root ~/work, --repos, ./.commit-chronicle,\n" +
			"  or ~/.config/commit-chronicle/{repos,roots}")
	}

	// Validate + de-duplicate (preserve discovery order).
	seen := make(map[string]bool)
	var valid []string
	for _, c := range candidates {
		p := expand(c)
		if seen[p] {
			continue
		}
		seen[p] = true
		if isGitRepo(p) {
			valid = append(valid, p)
		} else {
			fmt.Fprintf(os.Stderr, "⚠️  skipping (not a git repo): %s\n", p)
		}
	}
	if len(valid) == 0 {
		return nil, fmt.Errorf("no valid git repositories to scan")
	}
	return valid, nil
}

// discoverRepos walks root and returns every git repo beneath it (not
// descending into a repo once found, nor into heavy/build directories).
func discoverRepos(root string) []string {
	var repos []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil //nolint:nilerr // unreadable dirs are skipped, not fatal
		}
		base := filepath.Base(path)
		if path != root && (skipDirs[base] || strings.HasPrefix(base, ".")) {
			return fs.SkipDir
		}
		if isGitRepo(path) {
			repos = append(repos, path)
			// Stop descending into a nested repo (submodules etc.), but if the
			// root dir is itself a repo, keep going so sibling repos beneath it
			// (e.g. ~/work/forked/{a,b,c}) are still discovered.
			if path != root {
				return fs.SkipDir
			}
		}
		return nil
	})
	sort.Strings(repos)
	return repos
}

func xdgPath(name string) string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "commit-chronicle", name)
}

func linesFromFile(path string) ([]string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()

	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out, len(out) > 0
}

func expand(p string) string {
	p = strings.TrimSpace(p)
	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		p = filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}
	return p
}

func isGitRepo(path string) bool {
	if fi, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		_ = fi
		return true
	}
	// Fall back to git for worktrees / unusual layouts.
	out, err := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree").Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}
