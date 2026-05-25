// Package collect gathers worklog items from git history and GitHub:
// commits (by author and from authored PRs), authored PRs, and reviewed PRs.
package collect

import (
	"sort"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

// Options controls what Gather collects.
type Options struct {
	Repos          []string
	Author         string // git author name to match
	User           string // GitHub login (enables PR/review collection)
	Range          model.Range
	IncludePRs     bool // commits-from-PRs and authored-PR entries
	IncludeReviews bool // reviewed-PR entries
}

// Progress is an optional callback for status messages (may be nil).
type Progress func(stage string, count int)

// HasGH reports whether the gh CLI is available (for the caller to decide
// whether PR/review collection is even possible).
func HasGH() bool { return hasGH() }

// Gather collects everything requested into a de-duplicated, date-sorted slice.
//
// Sources, in order, with later sources only adding items not already seen:
//  1. git commits authored by Author across all refs
//  2. commits on the user's authored PRs        (if IncludePRs && gh)
//  3. the user's authored PRs as entries        (if IncludePRs && gh)
//  4. the user's reviewed PRs as entries        (if IncludeReviews && gh)
func Gather(o Options, p Progress) ([]model.Item, error) {
	report := func(stage string, n int) {
		if p != nil {
			p(stage, n)
		}
	}

	items := gitCommits(o.Repos, o.Author, o.Range)
	report("git commits", len(items))

	useGH := o.User != "" && hasGH()
	if useGH && o.IncludePRs {
		pc := prCommits(o.Repos, o.User, o.Range)
		report("PR commits", len(pc))
		items = append(items, pc...)

		ap := authoredPRs(o.Repos, o.User, o.Range)
		report("authored PRs", len(ap))
		items = append(items, ap...)
	}
	if useGH && o.IncludeReviews {
		rp := reviewedPRs(o.Repos, o.User, o.Range)
		report("reviewed PRs", len(rp))
		items = append(items, rp...)
	}

	return dedupeSort(items), nil
}

// dedupeSort removes duplicates by Item.ID and orders oldest→newest, then by
// kind (commits, PRs, reviews), then repo.
func dedupeSort(in []model.Item) []model.Item {
	seen := make(map[string]struct{}, len(in))
	out := in[:0]
	for _, it := range in {
		id := it.ID()
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, it)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Date != out[j].Date {
			return out[i].Date < out[j].Date
		}
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].RepoName < out[j].RepoName
	})
	return out
}
