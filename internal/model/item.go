// Package model defines the domain types shared across commit-chronicle:
// the unified worklog Item and the date Range.
package model

import "fmt"

// Kind distinguishes the sources that feed a worklog.
type Kind int

const (
	KindCommit Kind = iota // a git commit (history or PR head)
	KindPR                 // a pull request the user authored
	KindReview             // a pull request the user reviewed
)

// Item is one selectable, renderable entry in a worklog. A single struct
// covers commits and PRs so the picker and renderers stay uniform.
type Item struct {
	Kind     Kind
	Date     string // YYYY-MM-DD used for grouping and sorting
	RepoName string
	RepoPath string
	URL      string // canonical link (commit or PR), may be empty

	// Commit-only
	Hash      string
	ShortHash string

	// PR/Review-only
	Number int
	State  string // OPEN | MERGED | CLOSED

	// Common payload: commit subject or PR title
	Title string
}

// ID is the de-duplication key. Commits dedupe by hash; PRs/reviews by
// kind+repo+number so an authored PR and a reviewed PR never collide.
func (i Item) ID() string {
	if i.Kind == KindCommit {
		return "c:" + i.Hash
	}
	return fmt.Sprintf("%d:%s#%d", i.Kind, i.RepoName, i.Number)
}

// Tag is the short label shown in the picker.
func (i Item) Tag() string {
	switch i.Kind {
	case KindPR:
		return "PR"
	case KindReview:
		return "review"
	default:
		return "commit"
	}
}

// Ref is the compact identifier (short hash or #number).
func (i Item) Ref() string {
	if i.Kind == KindCommit {
		return i.ShortHash
	}
	return fmt.Sprintf("#%d", i.Number)
}
