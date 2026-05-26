package collect

import (
	"testing"

	"github.com/ashishxcode/commit-chronicle/internal/model"
)

func TestSlugFromURL(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"git@github.com:ashishxcode/cx-saas-dashboard.git", "ashishxcode/cx-saas-dashboard"},
		{"https://github.com/CultureX-art/cx-saas-dashboard.git", "CultureX-art/cx-saas-dashboard"},
		{"https://github.com/owner/repo", "owner/repo"},
		{"  git@github.com:owner/repo.git\n", "owner/repo"},
		{"git@gitlab.com:owner/repo.git", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := slugFromURL(tt.in); got != tt.want {
			t.Errorf("slugFromURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestUpdatedSinceQualifier(t *testing.T) {
	tests := []struct {
		name string
		r    model.Range
		want string
	}{
		{"single day uses open lower bound", model.Range{Since: "2026-05-25", Until: "2026-05-26"}, " updated:>=2026-05-25"},
		{"open-ended range", model.Range{Since: "2026-05-01"}, " updated:>=2026-05-01"},
		{"relative since is empty", model.Range{Since: "7 days ago"}, ""},
		{"empty since is empty", model.Range{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updatedSinceQualifier(tt.r); got != tt.want {
				t.Errorf("updatedSinceQualifier(%+v) = %q, want %q", tt.r, got, tt.want)
			}
		})
	}
}

func TestCreatedQualifier(t *testing.T) {
	tests := []struct {
		name string
		r    model.Range
		want string
	}{
		// Until is exclusive (next-day midnight); the inclusive upper bound is Until-1.
		{"single day is inclusive both ends", model.Range{Since: "2026-05-25", Until: "2026-05-26"}, " created:2026-05-25..2026-05-25"},
		{"multi-day window", model.Range{Since: "2026-05-01", Until: "2026-06-01"}, " created:2026-05-01..2026-05-31"},
		{"open-ended range", model.Range{Since: "2026-05-01"}, " created:>=2026-05-01"},
		{"relative since is empty", model.Range{Since: "7 days ago"}, ""},
		{"empty since is empty", model.Range{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createdQualifier(tt.r); got != tt.want {
				t.Errorf("createdQualifier(%+v) = %q, want %q", tt.r, got, tt.want)
			}
		})
	}
}
