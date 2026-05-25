package collect

import (
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

const fieldSep = "\x1f" // ASCII unit separator

// anchorMidnight pins a bare YYYY-MM-DD to local midnight. git's --since/--until
// parse a bare date via approxidate, which fills the missing time with the
// *current* time of day — skewing every window by the wall clock. Relative
// strings (e.g. "7 days ago") are passed through untouched.
func anchorMidnight(when string) string {
	if _, err := time.Parse("2006-01-02", when); err == nil {
		return when + " 00:00:00"
	}
	return when
}

// gitCommits returns de-duplicated commits authored by `author` in the range,
// across all refs in each repo.
func gitCommits(repos []string, author string, r model.Range) []model.Item {
	var items []model.Item
	for _, repo := range repos {
		name := filepath.Base(repo)
		base := originURL(repo)

		args := []string{
			"-C", repo, "log", "--all", "--no-merges",
			"--author=" + author, "--regexp-ignore-case",
			"--date=short",
			"--pretty=format:%h" + fieldSep + "%H" + fieldSep + "%ad" + fieldSep + "%s",
		}
		if r.Since != "" {
			args = append(args, "--since="+anchorMidnight(r.Since))
		}
		if r.Until != "" {
			args = append(args, "--until="+anchorMidnight(r.Until))
		}

		out, err := exec.Command("git", args...).Output()
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(out), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			p := strings.SplitN(line, fieldSep, 4)
			if len(p) != 4 {
				continue
			}
			url := ""
			if base != "" {
				url = base + "/commit/" + p[1]
			}
			items = append(items, model.Item{
				Kind:      model.KindCommit,
				Date:      p[2],
				RepoName:  name,
				RepoPath:  repo,
				URL:       url,
				Hash:      p[1],
				ShortHash: p[0],
				Title:     p[3],
			})
		}
	}
	return items
}

// Preview returns `git show --stat` for a commit item (for the picker pane).
func Preview(it model.Item) string {
	if it.Kind != model.KindCommit || it.Hash == "" {
		return previewPR(it)
	}
	out, err := exec.Command("git", "-C", it.RepoPath, "show", "--stat", "--no-color", it.Hash).Output()
	if err != nil {
		return "(commit not present locally — fetch the branch to see the diff)\n\n" + previewPR(it)
	}
	return string(out)
}

func previewPR(it model.Item) string {
	var b strings.Builder
	b.WriteString(it.Tag() + " " + it.Ref() + "\n")
	b.WriteString("repo:  " + it.RepoName + "\n")
	if it.State != "" {
		b.WriteString("state: " + it.State + "\n")
	}
	b.WriteString("date:  " + it.Date + "\n\n")
	b.WriteString(it.Title + "\n")
	if it.URL != "" {
		b.WriteString("\n" + it.URL + "\n")
	}
	return b.String()
}

// originURL returns the https GitHub URL for a repo's origin, or "".
func originURL(repoPath string) string {
	out, err := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	url := strings.TrimSpace(string(out))
	url = strings.TrimSuffix(url, ".git")
	url = strings.Replace(url, "git@github.com:", "https://github.com/", 1)
	return url
}
