package collect

import (
	"encoding/json"
	"os/exec"
	"sort"
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
		for _, slug := range repoSlugs(repo) {
			name := slug[strings.Index(slug, "/")+1:]
			base := "https://github.com/" + slug
			for _, n := range listPRNumbers(slug, "author:"+user+updatedSinceQualifier(r)) {
				for _, c := range viewPRCommits(slug, n) {
					day, ok := isoDay(c.AuthoredDate)
					if !ok || !inRange(day, since, until) {
						continue
					}
					if isNoiseSubject(c.MessageHeadline) {
						continue
					}
					short := c.OID
					if len(short) > 8 {
						short = short[:8]
					}
					items = append(items, model.Item{
						Kind:      model.KindCommit,
						Date:      day.Format("2006-01-02"),
						RepoName:  name,
						RepoPath:  repo,
						URL:       base + "/commit/" + c.OID,
						Hash:      c.OID,
						ShortHash: short,
						Title:     model.CleanText(c.MessageHeadline),
					})
				}
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
		for _, slug := range repoSlugs(repo) {
			name := slug[strings.Index(slug, "/")+1:]
			raw, err := exec.Command("gh", "pr", "list", "-R", slug,
				"--search", "author:"+user+createdQualifier(r), "--state", "all", "--limit", "200",
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
					Title:    model.CleanText(pr.Title),
				})
			}
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
		for _, slug := range repoSlugs(repo) {
			name := slug[strings.Index(slug, "/")+1:]

			// List PRs the user reviewed (numbers + meta only — cheap).
			raw, err := exec.Command("gh", "pr", "list", "-R", slug,
				"--search", "reviewed-by:"+user+updatedSinceQualifier(r), "--state", "all", "--limit", "200",
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
				day, verdict, ok := reviewDayInRange(slug, pr.Number, user, since, until)
				if !ok {
					continue
				}
				items = append(items, model.Item{
					Kind:        model.KindReview,
					Date:        day,
					RepoName:    name,
					RepoPath:    repo,
					URL:         pr.URL,
					Number:      pr.Number,
					State:       pr.State,
					ReviewState: verdict,
					Title:       model.CleanText(pr.Title),
				})
			}
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

// reviewDayInRange returns the date and verdict (APPROVED, CHANGES_REQUESTED,
// COMMENTED, …) of the user's earliest in-range review on a PR. Any review
// state counts — a "changes requested" review is as much a review as an
// approval — and the verdict reflects the dated (earliest in-range) review.
func reviewDayInRange(slug string, number int, user string, since, until time.Time) (string, string, bool) {
	raw, err := exec.Command("gh", "pr", "view", strconv.Itoa(number),
		"-R", slug, "--json", "reviews").Output()
	if err != nil {
		return "", "", false
	}
	var v struct {
		Reviews []ghReview `json:"reviews"`
	}
	if json.Unmarshal(raw, &v) != nil {
		return "", "", false
	}
	best, verdict := "", ""
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
			best, verdict = d, rv.State
		}
	}
	return best, verdict, best != ""
}

// slugFromURL extracts a GitHub "owner/repo" from a remote URL, or "".
func slugFromURL(url string) string {
	url = strings.TrimSuffix(strings.TrimSpace(url), ".git")
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

// repoSlugs returns the distinct GitHub "owner/repo" slugs across all of a
// repo's remotes, origin first. Fork workflows push to a personal "origin" but
// open PRs and submit reviews against the "upstream" parent, so PR and review
// discovery must consider every remote — querying origin alone misses them.
func repoSlugs(repoPath string) []string {
	out, err := exec.Command("git", "-C", repoPath, "remote").Output()
	if err != nil {
		return nil
	}
	remotes := strings.Fields(string(out))
	// Query origin first so its repo name drives display naming.
	sort.SliceStable(remotes, func(i, j int) bool {
		return remotes[i] == "origin" && remotes[j] != "origin"
	})
	var slugs []string
	seen := map[string]struct{}{}
	for _, name := range remotes {
		u, err := exec.Command("git", "-C", repoPath, "remote", "get-url", name).Output()
		if err != nil {
			continue
		}
		s := slugFromURL(string(u))
		if s == "" {
			continue
		}
		if _, dup := seen[s]; dup {
			continue
		}
		seen[s] = struct{}{}
		slugs = append(slugs, s)
	}
	return slugs
}

// --- date helpers ---------------------------------------------------------

// updatedSinceQualifier bounds a PR search to those touched on/after Since, to
// trim fan-out without dropping in-range items. Any PR carrying an in-range
// commit or review has updatedAt >= Since, so the lower bound is safe. We
// deliberately omit an upper bound: a PR reviewed (or committed to) within the
// window can be updated again afterward, and "updated:<=Until" would then
// wrongly hide it before the precise per-item date filtering runs. Empty when
// Since isn't a concrete date (e.g. a relative "7 days ago").
func updatedSinceQualifier(r model.Range) string {
	if parseDay(r.Since).IsZero() {
		return ""
	}
	return " updated:>=" + r.Since
}

// createdQualifier bounds an authored-PR search by creation date — exactly the
// field authoredPRs dates its items by, so both bounds are precise. Until is
// exclusive (next-day midnight); GitHub's created:A..B is inclusive, so the
// upper bound is the day before Until. Empty for a non-concrete Since.
func createdQualifier(r model.Range) string {
	since, until := parseDay(r.Since), parseDay(r.Until)
	if since.IsZero() {
		return ""
	}
	if !until.IsZero() {
		hi := until.AddDate(0, 0, -1).Format("2006-01-02")
		return " created:" + r.Since + ".." + hi
	}
	return " created:>=" + r.Since
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
