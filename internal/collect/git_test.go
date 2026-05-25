package collect

import "testing"

func TestAnchorMidnight(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"bare ISO date gets midnight", "2026-05-24", "2026-05-24 00:00:00"},
		{"month rollover date", "2026-02-01", "2026-02-01 00:00:00"},
		{"relative string untouched", "7 days ago", "7 days ago"},
		{"yesterday keyword untouched", "yesterday", "yesterday"},
		{"empty untouched", "", ""},
		{"date with time untouched", "2026-05-24 12:00:00", "2026-05-24 12:00:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := anchorMidnight(tt.in); got != tt.want {
				t.Errorf("anchorMidnight(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
