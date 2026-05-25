package collect

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

// maxPRs bounds how many PRs we inspect per repo, to keep API calls sane.
const maxPRs = 150

// hasGH reports whether the gh CLI is available.
func hasGH() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// --- gh JSON shapes -------------------------------------------------------

type ghNum struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	State     string `json:"state"`
	URL       string `json:"url"`
	CreatedAt string `json:"createdAt"`
}

type ghCommit struct {
	OID             string `json:"oid"`
	MessageHeadline string `json:"messageHeadline"`
	AuthoredDate    string `json:"authoredDate"`
}

type ghReview struct {
	State       string `json:"state"`
	SubmittedAt string `json:"submittedAt"`
	Author      struct {
		Login string `json:"login"`
	} `json:"author"`
}

// prCommits returns commits on the user's authored PRs (KindCommit), so commits
// the plain author filter misses (different email, unfetched branch) are kept.
func prCommits(repos []string, user string, r model.Range) []model.Item {
	since, until := parseDay(r.Since), parseDay(r.Until)
	var items []model.Item
	for _, repo := range repos {
		slug := repoSlug(repo)
		if slug == "" {
			continue
		}
		name := slug[strings.Index(slug, "/")+1:]
		base := originURL(repo)
		for _, n := range listPRNumbers(slug, "author:"+user+dateQualifier(r)) {
			for _, c := range viewPRCommits(slug, n) {
				day, ok := isoDay(c.AuthoredDate)
				if !ok || !inRange(day, since, until) {
					continue
				}
				short := c.OID
				if len(short) > 8 {
					short = short[:8]
				}
				url := ""
				if base != "" {
					url = base + "/commit/" + c.OID
				}
				items = append(items, model.Item{
					Kind:      model.KindCommit,
					Date:      day.Format("2006-01-02"),
					RepoName:  name,
					RepoPath:  repo,
					URL:       url,
					Hash:      c.OID,
					ShortHash: short,
					Title:     c.MessageHeadline,
				})
			}
		}
	}
	return items
}

// authoredPRs returns PRs the user opened, created within the range (KindPR).
func authoredPRs(repos []string, user string, r model.Range) []model.Item {
	since, until := parseDay(r.Since), parseDay(r.Until)
	var items []model.Item
	for _, repo := range repos {
		slug := repoSlug(repo)
		if slug == "" {
			continue
		}
		name := slug[strings.Index(slug, "/")+1:]
		raw, err := exec.Command("gh", "pr", "list", "-R", slug,
			"--search", "author:"+user+dateQualifier(r), "--state", "all", "--limit", "200",
			"--json", "number,title,state,url,createdAt").Output()
		if err != nil {
			continue
		}
		var prs []ghNum
		if json.Unmarshal(raw, &prs) != nil {
			continue
		}
		for _, pr := range prs {
			day, ok := isoDay(pr.CreatedAt)
			if !ok || !inRange(day, since, until) {
				continue
			}
			items = append(items, model.Item{
				Kind:     model.KindPR,
				Date:     day.Format("2006-01-02"),
				RepoName: name,
				RepoPath: repo,
				URL:      pr.URL,
				Number:   pr.Number,
				State:    pr.State,
				Title:    pr.Title,
			})
		}
	}
	return items
}

// reviewedPRs returns PRs the user reviewed within the range (KindReview).
// Dated by the user's review submission time, one review per PR (earliest in range).
func reviewedPRs(repos []string, user string, r model.Range) []model.Item {
	since, until := parseDay(r.Since), parseDay(r.Until)
	var items []model.Item
	for _, repo := range repos {
		slug := repoSlug(repo)
		if slug == "" {
			continue
		}
		name := slug[strings.Index(slug, "/")+1:]

		// List PRs the user reviewed (numbers + meta only — cheap).
		raw, err := exec.Command("gh", "pr", "list", "-R", slug,
			"--search", "reviewed-by:"+user+dateQualifier(r), "--state", "all", "--limit", "200",
			"--json", "number,title,state,url").Output()
		if err != nil {
			continue
		}
		var prs []ghNum
		if json.Unmarshal(raw, &prs) != nil {
			continue
		}
		if len(prs) > maxPRs {
			prs = prs[:maxPRs]
		}
		for _, pr := range prs {
			day, ok := reviewDayInRange(slug, pr.Number, user, since, until)
			if !ok {
				continue
			}
			items = append(items, model.Item{
				Kind:     model.KindReview,
				Date:     day,
				RepoName: name,
				RepoPath: repo,
				URL:      pr.URL,
				Number:   pr.Number,
				State:    pr.State,
				Title:    pr.Title,
			})
		}
	}
	return items
}

// --- gh call helpers ------------------------------------------------------

func listPRNumbers(slug, search string) []int {
	raw, err := exec.Command("gh", "pr", "list", "-R", slug,
		"--search", search, "--state", "all", "--limit", "200",
		"--json", "number").Output()
	if err != nil {
		return nil
	}
	var nums []ghNum
	if json.Unmarshal(raw, &nums) != nil {
		return nil
	}
	if len(nums) > maxPRs {
		nums = nums[:maxPRs]
	}
	out := make([]int, len(nums))
	for i, n := range nums {
		out[i] = n.Number
	}
	return out
}

func viewPRCommits(slug string, number int) []ghCommit {
	raw, err := exec.Command("gh", "pr", "view", strconv.Itoa(number),
		"-R", slug, "--json", "commits").Output()
	if err != nil {
		return nil
	}
	var v struct {
		Commits []ghCommit `json:"commits"`
	}
	if json.Unmarshal(raw, &v) != nil {
		return nil
	}
	return v.Commits
}

// reviewDayInRange returns the date of the user's earliest in-range review on a PR.
func reviewDayInRange(slug string, number int, user string, since, until time.Time) (string, bool) {
	raw, err := exec.Command("gh", "pr", "view", strconv.Itoa(number),
		"-R", slug, "--json", "reviews").Output()
	if err != nil {
		return "", false
	}
	var v struct {
		Reviews []ghReview `json:"reviews"`
	}
	if json.Unmarshal(raw, &v) != nil {
		return "", false
	}
	best := ""
	for _, rv := range v.Reviews {
		if !strings.EqualFold(rv.Author.Login, user) {
			continue
		}
		day, ok := isoDay(rv.SubmittedAt)
		if !ok || !inRange(day, since, until) {
			continue
		}
		d := day.Format("2006-01-02")
		if best == "" || d < best {
			best = d
		}
	}
	return best, best != ""
}

// repoSlug returns "owner/repo" for a repo's GitHub origin, or "".
func repoSlug(repoPath string) string {
	out, err := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	url := strings.TrimSpace(string(out))
	url = strings.TrimSuffix(url, ".git")
	switch {
	case strings.HasPrefix(url, "git@github.com:"):
		url = strings.TrimPrefix(url, "git@github.com:")
	case strings.Contains(url, "github.com/"):
		url = url[strings.Index(url, "github.com/")+len("github.com/"):]
	default:
		return ""
	}
	if parts := strings.SplitN(url, "/", 2); len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0] + "/" + parts[1]
	}
	return ""
}

// --- date helpers ---------------------------------------------------------

// dateQualifier returns a GitHub search suffix bounding PRs to the window by
// last-updated date, so we don't fan out gh calls over out-of-range PRs. Empty
// when Since isn't a concrete date (e.g. a relative "7 days ago").
func dateQualifier(r model.Range) string {
	if parseDay(r.Since).IsZero() {
		return ""
	}
	if !parseDay(r.Until).IsZero() {
		return " updated:" + r.Since + ".." + r.Until
	}
	return " updated:>=" + r.Since
}

func parseDay(s string) time.Time {
	if s == "" || s == "now" {
		return time.Time{}
	}
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func isoDay(s string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false
	}
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), true
}

func inRange(day, since, until time.Time) bool {
	if !since.IsZero() && day.Before(since) {
		return false
	}
	if !until.IsZero() && !day.Before(until) { // until exclusive
		return false
	}
	return true
}
