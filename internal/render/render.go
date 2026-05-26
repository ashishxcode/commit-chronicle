// Package render turns selected worklog items into markdown or json.
package render

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

// Meta is the worklog header information.
type Meta struct {
	Author     string
	RangeLabel string
}

// reviewVerdict maps a review state to a display icon and label.
func reviewVerdict(state string) (icon, label string) {
	switch state {
	case "APPROVED":
		return "✅", "approved"
	case "CHANGES_REQUESTED":
		return "🔴", "changes requested"
	case "COMMENTED":
		return "💬", "commented"
	case "DISMISSED":
		return "⊘", "dismissed"
	default:
		return "•", "reviewed"
	}
}

func stateIcon(state string) string {
	switch state {
	case "MERGED":
		return "✓"
	case "OPEN":
		return "○"
	case "CLOSED":
		return "✕"
	default:
		return "•"
	}
}

// Markdown renders items grouped by date, then by kind within each day, so a
// reader sees a clean "commits / pull requests / reviews" breakdown per day.
func Markdown(items []model.Item, m Meta) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 📓 Worklog — %s\n\n", m.RangeLabel)
	fmt.Fprintf(&b, "**%s** · generated %s  \n", m.Author, time.Now().Format("2006-01-02 15:04"))
	fmt.Fprintf(&b, "%s\n", counts(items))

	// Walk the (already date-sorted, then kind-sorted) items, opening a new day
	// section on date change and a kind subsection on kind change within a day.
	curDate, curKind := "", model.Kind(-1)
	for _, it := range items {
		if it.Date != curDate {
			curDate, curKind = it.Date, model.Kind(-1)
			fmt.Fprintf(&b, "\n## %s\n", it.Date)
		}
		if it.Kind != curKind {
			curKind = it.Kind
			fmt.Fprintf(&b, "\n### %s\n\n", kindHeading(it.Kind))
		}
		b.WriteString(line(it))
	}
	b.WriteString("\n")
	return b.String()
}

func kindHeading(k model.Kind) string {
	switch k {
	case model.KindPR:
		return "Pull requests"
	case model.KindReview:
		return "Reviews"
	default:
		return "Commits"
	}
}

func line(it model.Item) string {
	link := func(text string) string {
		if it.URL != "" {
			return fmt.Sprintf("[%s](%s)", text, it.URL)
		}
		return text
	}
	switch it.Kind {
	case model.KindPR:
		return fmt.Sprintf("- %s **%s** %s — %s · %s\n",
			stateIcon(it.State), link(it.Ref()), it.Title, strings.ToLower(it.State), it.RepoName)
	case model.KindReview:
		icon, verdict := reviewVerdict(it.ReviewState)
		return fmt.Sprintf("- %s **%s** %s — %s (PR %s) · %s\n",
			icon, link(it.Ref()), it.Title, verdict, strings.ToLower(it.State), it.RepoName)
	default: // commit
		return fmt.Sprintf("- %s  _(%s)_\n", it.Title,
			link(it.RepoName+"@"+it.ShortHash))
	}
}

func counts(items []model.Item) string {
	var c, p, r int
	for _, it := range items {
		switch it.Kind {
		case model.KindCommit:
			c++
		case model.KindPR:
			p++
		case model.KindReview:
			r++
		}
	}
	return fmt.Sprintf("%d commits · %d PRs · %d reviews · %d total", c, p, r, len(items))
}

type jsonItem struct {
	Kind        string `json:"kind"`
	Date        string `json:"date"`
	Repo        string `json:"repo"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Hash        string `json:"hash,omitempty"`
	Number      int    `json:"number,omitempty"`
	State       string `json:"state,omitempty"`
	ReviewState string `json:"reviewState,omitempty"`
}

// JSON renders items as a JSON array.
func JSON(items []model.Item, _ Meta) string {
	out := make([]jsonItem, 0, len(items))
	for _, it := range items {
		out = append(out, jsonItem{
			Kind: it.Tag(), Date: it.Date, Repo: it.RepoName,
			Title: it.Title, URL: it.URL, Hash: it.Hash,
			Number: it.Number, State: it.State, ReviewState: it.ReviewState,
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "[]\n"
	}
	return string(data) + "\n"
}
