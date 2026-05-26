package render

import (
	"strings"
	"testing"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

func TestReviewVerdict(t *testing.T) {
	tests := []struct {
		state, wantLabel string
	}{
		{"APPROVED", "approved"},
		{"CHANGES_REQUESTED", "changes requested"},
		{"COMMENTED", "commented"},
		{"DISMISSED", "dismissed"},
		{"", "reviewed"},
		{"WEIRD_STATE", "reviewed"},
	}
	for _, tt := range tests {
		if _, got := reviewVerdict(tt.state); got != tt.wantLabel {
			t.Errorf("reviewVerdict(%q) label = %q, want %q", tt.state, got, tt.wantLabel)
		}
	}
}

// A changes-requested review must render with its verdict (not the PR's
// merge state) so it is clearly included and distinguishable from approvals.
func TestReviewLineShowsVerdict(t *testing.T) {
	it := model.Item{
		Kind:        model.KindReview,
		Date:        "2026-05-22",
		RepoName:    "cx-saas-dashboard",
		Number:      2319,
		State:       "OPEN",
		ReviewState: "CHANGES_REQUESTED",
		Title:       "Fix: bucket group image removal",
	}
	got := line(it)
	if !strings.Contains(got, "changes requested") {
		t.Errorf("review line missing verdict: %q", got)
	}
	if !strings.Contains(got, "PR open") {
		t.Errorf("review line missing PR state: %q", got)
	}
}
