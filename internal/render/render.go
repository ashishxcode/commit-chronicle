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

// Markdown renders items grouped by date. Commits, PRs and reviews each get a
// distinct, link-bearing line so nothing is ambiguous.
func Markdown(items []model.Item, m Meta) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 📓 Worklog — %s\n\n", m.RangeLabel)
	fmt.Fprintf(&b, "> **Author:** %s  \n", m.Author)
	fmt.Fprintf(&b, "> **Generated:** %s  \n", time.Now().Format("2006-01-02 15:04"))
	fmt.Fprintf(&b, "> **Entries:** %d  (%s)\n\n", len(items), counts(items))
	b.WriteString("<!-- Edit freely. Each '- ' line is a worklog entry. -->\n")

	cur := ""
	for _, it := range items {
		if it.Date != cur {
			cur = it.Date
			fmt.Fprintf(&b, "\n## %s\n\n", it.Date)
		}
		b.WriteString(line(it))
	}
	b.WriteString("\n")
	return b.String()
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
		return fmt.Sprintf("- **PR %s** %s %s  _(%s)_ — %s\n",
			link(it.Ref()), stateIcon(it.State), it.Title, it.State, it.RepoName)
	case model.KindReview:
		return fmt.Sprintf("- **Reviewed PR %s** %s %s  _(%s)_ — %s\n",
			link(it.Ref()), stateIcon(it.State), it.Title, it.State, it.RepoName)
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
	return fmt.Sprintf("%d commits, %d PRs, %d reviews", c, p, r)
}

type jsonItem struct {
	Kind   string `json:"kind"`
	Date   string `json:"date"`
	Repo   string `json:"repo"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Hash   string `json:"hash,omitempty"`
	Number int    `json:"number,omitempty"`
	State  string `json:"state,omitempty"`
}

// JSON renders items as a JSON array.
func JSON(items []model.Item, _ Meta) string {
	out := make([]jsonItem, 0, len(items))
	for _, it := range items {
		out = append(out, jsonItem{
			Kind: it.Tag(), Date: it.Date, Repo: it.RepoName,
			Title: it.Title, URL: it.URL, Hash: it.Hash,
			Number: it.Number, State: it.State,
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "[]\n"
	}
	return string(data) + "\n"
}
