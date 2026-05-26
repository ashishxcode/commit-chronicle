package collect

import "testing"

func TestIsNoiseSubject(t *testing.T) {
	noise := []string{
		"Merge branch 'development' into feature",
		"Merge remote-tracking branch 'origin/main'",
		"Merge pull request #42 from foo/bar",
		"index on main: abc123 do thing",
		"WIP on feature: abc123 in progress",
	}
	for _, s := range noise {
		if !isNoiseSubject(s) {
			t.Errorf("isNoiseSubject(%q) = false, want true", s)
		}
	}
	real := []string{
		"fix: redirect external LinkTracker URLs",
		"feat: add createMutation factory",
		"Mergesort optimization", // not a merge commit
		"index page layout fix",  // not a stash
	}
	for _, s := range real {
		if isNoiseSubject(s) {
			t.Errorf("isNoiseSubject(%q) = true, want false", s)
		}
	}
}

func TestAnchorMidnight(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"bare ISO date gets midnight", "2026-05-24", "2026-05-24 00:00:00"},
		{"month rollover date", "2026-02-01", "2026-02-01 00:00:00"},
		{"relative string anchored", "7 days ago", "7 days ago 00:00:00"},
		{"today keyword anchored", "today", "today 00:00:00"},
		{"yesterday keyword anchored", "yesterday", "yesterday 00:00:00"},
		{"empty untouched", "", ""},
		{"date with time untouched", "2026-05-24 12:00:00", "2026-05-24 12:00:00"},
		{"midnight keyword untouched", "midnight", "midnight"},
		{"noon keyword untouched", "noon", "noon"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := anchorMidnight(tt.in); got != tt.want {
				t.Errorf("anchorMidnight(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
