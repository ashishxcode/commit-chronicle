// Package model defines the domain types shared across commit-chronicle:
// the unified worklog Item and the date Range.
package model

import (
	"fmt"
	"regexp"
	"strings"
)

// ansiEscape matches terminal escape sequences: CSI (ESC [ … final byte),
// OSC (ESC ] … terminated by BEL or ST), and simple two-byte escapes. These
// are stripped whole so no visible residue (e.g. "[31m") is left behind.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]` + // CSI
	`|\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)?` + // OSC … BEL or ST
	`|\x1b[@-Z\\-_]`) // other single-char escapes

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

	// Review-only: your verdict on the PR (APPROVED | CHANGES_REQUESTED |
	// COMMENTED | DISMISSED), from the review we date this entry by.
	ReviewState string

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

// CleanText strips control and escape characters from externally-sourced text
// (commit subjects, PR titles). Without this, a crafted commit message or PR
// title could inject terminal escape sequences — moving the cursor, hiding
// output, or rewriting the screen — when shown in the picker, the preview pane,
// or a worklog printed to the terminal. Tabs become spaces; all C0/C1 control
// characters (including ESC, CR and LF) are dropped; printable text is kept.
func CleanText(s string) string {
	s = ansiEscape.ReplaceAllString(s, "")
	cleaned := strings.Map(func(r rune) rune {
		switch {
		case r == '\t':
			return ' '
		case r < 0x20, r == 0x7f, r >= 0x80 && r <= 0x9f:
			return -1
		default:
			return r
		}
	}, s)
	return strings.TrimSpace(cleaned)
}
