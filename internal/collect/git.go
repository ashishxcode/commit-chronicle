package collect

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

const fieldSep = "\x1f" // ASCII unit separator

// isNoiseSubject reports whether a commit subject is mechanical noise that
// shouldn't appear in a worklog: merge commits and git-stash entries. (Plain
// git history is already filtered with --no-merges; this also covers commits
// pulled in from PRs, and stash refs that slip in via `git log --all`.)
func isNoiseSubject(s string) bool {
	s = strings.TrimSpace(s)
	for _, p := range []string{
		"Merge branch ", "Merge remote-tracking ", "Merge pull request ",
		"index on ", "WIP on ",
	} {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// anchorMidnight pins a date or relative day-phrase to local midnight. git's
// --since/--until parse their argument via approxidate, which fills a missing
// time-of-day with the *current* wall-clock time — so a bare YYYY-MM-DD,
// "today", "yesterday", or "7 days ago" all silently skew the window by the
// time of day (e.g. "--since today" excludes everything committed earlier
// today). Appending an explicit 00:00:00 anchors to the start of the day,
// matching the --date=short granularity we report. Strings that already carry
// a time component (a ":" or the "midnight"/"noon" keywords) are left as-is.
func anchorMidnight(when string) string {
	if when == "" {
		return when
	}
	w := strings.ToLower(when)
	if strings.Contains(w, ":") || strings.Contains(w, "midnight") || strings.Contains(w, "noon") {
		return when
	}
	return when + " 00:00:00"
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
			if isNoiseSubject(p[3]) {
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
				Title:     model.CleanText(p[3]),
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
