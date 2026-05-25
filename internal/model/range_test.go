package model

import (
	"testing"
	"time"
)

func TestFromDates(t *testing.T) {
	tests := []struct {
		name      string
		from, to  string
		wantSince string
		wantUntil string
		wantLabel string
		wantErr   bool
	}{
		{
			name: "single day makes until exclusive next day",
			from: "2026-05-24", to: "2026-05-24",
			wantSince: "2026-05-24", wantUntil: "2026-05-25", wantLabel: "2026-05-24",
		},
		{
			name: "multi-day range",
			from: "2026-05-01", to: "2026-05-03",
			wantSince: "2026-05-01", wantUntil: "2026-05-04", wantLabel: "2026-05-01 → 2026-05-03",
		},
		{
			name: "until rolls over month boundary",
			from: "2026-01-31", to: "2026-01-31",
			wantSince: "2026-01-31", wantUntil: "2026-02-01", wantLabel: "2026-01-31",
		},
		{name: "bad from", from: "nope", to: "2026-05-24", wantErr: true},
		{name: "bad to", from: "2026-05-24", to: "nope", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := FromDates(tt.from, tt.to)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", r)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Since != tt.wantSince || r.Until != tt.wantUntil || r.Label != tt.wantLabel {
				t.Errorf("got {Since:%q Until:%q Label:%q}, want {Since:%q Until:%q Label:%q}",
					r.Since, r.Until, r.Label, tt.wantSince, tt.wantUntil, tt.wantLabel)
			}
		})
	}
}

func TestFromDatesEmptyToDefaultsToday(t *testing.T) {
	r, err := FromDates("2026-05-01", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantUntil := time.Now().AddDate(0, 0, 1).Format(isoDate)
	if r.Until != wantUntil {
		t.Errorf("Until = %q, want %q (today + 1, exclusive)", r.Until, wantUntil)
	}
}

func TestFromMonth(t *testing.T) {
	r, err := FromMonth("2026-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Since != "2026-02-01" || r.Until != "2026-03-01" {
		t.Errorf("got {Since:%q Until:%q}, want {2026-02-01 2026-03-01}", r.Since, r.Until)
	}
	if r.Label != "February 2026" {
		t.Errorf("Label = %q, want %q", r.Label, "February 2026")
	}
	if _, err := FromMonth("bad"); err == nil {
		t.Error("expected error for bad month")
	}
}

// TestPresetBoundaries guards the regression behind the date-skew bug: each
// preset must produce a half-open [Since, Until) window aligned to whole days,
// independent of the current time of day.
func TestPresetBoundaries(t *testing.T) {
	now := time.Now()
	today := now.Format(isoDate)
	yesterday := now.AddDate(0, 0, -1).Format(isoDate)
	tomorrow := now.AddDate(0, 0, 1).Format(isoDate)

	tests := []struct {
		choice    string
		wantSince string
		wantUntil string
	}{
		{"Today", today, tomorrow},
		{"Yesterday", yesterday, today},
		{"Last 7 days", now.AddDate(0, 0, -6).Format(isoDate), tomorrow},
		{"Last 30 days", now.AddDate(0, 0, -29).Format(isoDate), tomorrow},
	}
	for _, tt := range tests {
		t.Run(tt.choice, func(t *testing.T) {
			r := Preset(tt.choice)
			if r.Since != tt.wantSince {
				t.Errorf("Since = %q, want %q", r.Since, tt.wantSince)
			}
			if r.Until != tt.wantUntil {
				t.Errorf("Until = %q, want %q", r.Until, tt.wantUntil)
			}
		})
	}
}

func TestPresetUnknownFallsBackToWeek(t *testing.T) {
	r := Preset("not a preset")
	want := Preset("Last 7 days")
	if r.Since != want.Since || r.Until != want.Until {
		t.Errorf("unknown preset = {%q,%q}, want last-7-days {%q,%q}",
			r.Since, r.Until, want.Since, want.Until)
	}
}
