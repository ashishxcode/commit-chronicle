package model

import "testing"

func TestCleanText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text unchanged", "fix: add column button", "fix: add column button"},
		{"strips ANSI color sequence", "title\x1b[31mRED\x1b[0m", "titleRED"},
		{"strips cursor-move escape whole", "\x1b[2Joops", "oops"},
		{"strips OSC title-set sequence", "x\x1b]0;pwned\x07y", "xy"},
		{"drops newlines and CR", "line1\r\nline2", "line1line2"},
		{"tab becomes space", "a\tb", "a b"},
		{"trims surrounding space", "  hi  ", "hi"},
		{"drops DEL and C1", "x\x7fy", "xy"},
		{"keeps unicode", "résumé ✓ 日本語", "résumé ✓ 日本語"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanText(tt.in); got != tt.want {
				t.Errorf("CleanText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
