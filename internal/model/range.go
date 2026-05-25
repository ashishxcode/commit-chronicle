package model

import (
	"fmt"
	"time"
)

const isoDate = "2006-01-02"

// Range is a reporting window. Until is exclusive (git semantics). Since/Until
// are YYYY-MM-DD where possible (a relative Since like "7 days ago" is allowed
// but disables PR date-filtering precision).
type Range struct {
	Since string
	Until string
	Label string
}

// FromDates builds an inclusive [from, to] range (to defaults to today).
func FromDates(from, to string) (Range, error) {
	fd, err := time.Parse(isoDate, from)
	if err != nil {
		return Range{}, fmt.Errorf("bad from-date %q (use YYYY-MM-DD)", from)
	}
	if to == "" {
		to = time.Now().Format(isoDate)
	}
	td, err := time.Parse(isoDate, to)
	if err != nil {
		return Range{}, fmt.Errorf("bad to-date %q (use YYYY-MM-DD)", to)
	}
	label := from + " → " + to
	if from == to {
		label = from
	}
	return Range{
		Since: fd.Format(isoDate),
		Until: td.AddDate(0, 0, 1).Format(isoDate), // make 'to' inclusive
		Label: label,
	}, nil
}

// FromMonth builds a whole-calendar-month range from YYYY-MM.
func FromMonth(month string) (Range, error) {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return Range{}, fmt.Errorf("bad month %q (use YYYY-MM)", month)
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	next := start.AddDate(0, 1, 0)
	return Range{
		Since: start.Format(isoDate),
		Until: next.Format(isoDate),
		Label: start.Format("January 2006"),
	}, nil
}

// Preset maps an interactive menu choice to a concrete range.
func Preset(choice string) Range {
	now := time.Now()
	today := now.Format(isoDate)
	switch choice {
	case "Today":
		r, _ := FromDates(today, today)
		return r
	case "Yesterday":
		y := now.AddDate(0, 0, -1).Format(isoDate)
		r, _ := FromDates(y, y)
		return r
	case "Last 7 days":
		r, _ := FromDates(now.AddDate(0, 0, -6).Format(isoDate), today)
		r.Label = "last 7 days"
		return r
	case "Last 30 days":
		r, _ := FromDates(now.AddDate(0, 0, -29).Format(isoDate), today)
		r.Label = "last 30 days"
		return r
	case "This month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		return Range{Since: start.Format(isoDate), Label: start.Format("January 2006")}
	case "This year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		return Range{Since: start.Format(isoDate), Label: fmt.Sprintf("%d", now.Year())}
	default:
		r, _ := FromDates(now.AddDate(0, 0, -6).Format(isoDate), today)
		r.Label = "last 7 days"
		return r
	}
}

// CustomRangeLabel is the menu entry that triggers the custom date selector.
const CustomRangeLabel = "Custom range…"

// PresetNames are the interactive range options, in display order.
var PresetNames = []string{"Today", "Yesterday", "Last 7 days", "Last 30 days", "This month", "This year", CustomRangeLabel}
