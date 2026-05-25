// Package tui implements the interactive screens (range, picker, editor).
package tui

import (
	"fmt"
	"strings"

	"github.com/ashishxcode/commit-chronicle/internal/collect"
	"github.com/ashishxcode/commit-chronicle/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cAccent  = lipgloss.Color("13")
	cDim     = lipgloss.Color("245")
	cSel     = lipgloss.Color("10")
	cBlue    = lipgloss.Color("12")
	titleSty = lipgloss.NewStyle().Bold(true).Foreground(cAccent)
	dimSty   = lipgloss.NewStyle().Foreground(cDim)
	selSty   = lipgloss.NewStyle().Foreground(cSel)
	curSty   = lipgloss.NewStyle().Bold(true).Foreground(cBlue)
	boxSty   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(cDim)
	tagSty   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

type pickerModel struct {
	items        []model.Item
	filtered     []int
	selected     map[string]bool // keyed by Item.ID()
	cursor       int
	top          int
	filter       textinput.Model
	filtering    bool
	preview      viewport.Model
	previewCache map[string]string
	width        int
	height       int
	rangeLabel   string
	author       string
	canceled     bool
	ready        bool
}

// Pick runs the interactive multi-select picker over commits and PRs.
func Pick(items []model.Item, rangeLabel, author string) ([]model.Item, bool, error) {
	ti := textinput.New()
	ti.Placeholder = "type to filter…"
	ti.Prompt = "filter ❯ "

	m := pickerModel{
		items:        items,
		selected:     make(map[string]bool),
		previewCache: make(map[string]string),
		filter:       ti,
		rangeLabel:   rangeLabel,
		author:       author,
	}
	m.applyFilter()

	res, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return nil, false, err
	}
	fm := res.(pickerModel)
	if fm.canceled {
		return nil, true, nil
	}
	var out []model.Item
	for _, it := range items {
		if fm.selected[it.ID()] {
			out = append(out, it)
		}
	}
	return out, false, nil
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m *pickerModel) applyFilter() {
	q := strings.ToLower(strings.TrimSpace(m.filter.Value()))
	tokens := strings.Fields(q)
	m.filtered = m.filtered[:0]
	for i, it := range m.items {
		hay := strings.ToLower(it.Tag() + " " + it.Date + " " + it.RepoName + " " + it.Ref() + " " + it.Title)
		ok := true
		for _, t := range tokens {
			if !strings.Contains(hay, t) {
				ok = false
				break
			}
		}
		if ok {
			m.filtered = append(m.filtered, i)
		}
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
	m.clampScroll()
}

func (m *pickerModel) listRows() int {
	r := m.height - 5
	if r < 3 {
		r = 3
	}
	return r
}

func (m *pickerModel) clampScroll() {
	rows := m.listRows()
	if m.cursor < m.top {
		m.top = m.cursor
	}
	if m.cursor >= m.top+rows {
		m.top = m.cursor - rows + 1
	}
	if m.top < 0 {
		m.top = 0
	}
}

func (m *pickerModel) updatePreview() {
	if len(m.filtered) == 0 {
		m.preview.SetContent(dimSty.Render("no items match the filter"))
		return
	}
	it := m.items[m.filtered[m.cursor]]
	body, ok := m.previewCache[it.ID()]
	if !ok {
		body = collect.Preview(it)
		m.previewCache[it.ID()] = body
	}
	m.preview.SetContent(body)
	m.preview.GotoTop()
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		listW := m.width/2 - 2
		prevW := m.width - listW - 4
		if !m.ready {
			m.preview = viewport.New(prevW, m.listRows())
			m.ready = true
		} else {
			m.preview.Width = prevW
			m.preview.Height = m.listRows()
		}
		m.filter.Width = m.width - 12
		m.clampScroll()
		m.updatePreview()
		return m, nil

	case tea.KeyMsg:
		if m.filtering {
			switch msg.String() {
			case "enter", "esc":
				m.filtering = false
				m.filter.Blur()
				return m, nil
			default:
				var cmd tea.Cmd
				m.filter, cmd = m.filter.Update(msg)
				m.applyFilter()
				m.updatePreview()
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.canceled = true
			return m, tea.Quit
		case "/":
			m.filtering = true
			m.filter.Focus()
			return m, textinput.Blink
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.clampScroll()
				m.updatePreview()
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.clampScroll()
				m.updatePreview()
			}
		case "pgup":
			m.preview.HalfViewUp()
		case "pgdown":
			m.preview.HalfViewDown()
		case " ", "tab":
			if len(m.filtered) > 0 {
				id := m.items[m.filtered[m.cursor]].ID()
				m.selected[id] = !m.selected[id]
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
					m.clampScroll()
					m.updatePreview()
				}
			}
		case "a":
			allSel := true
			for _, idx := range m.filtered {
				if !m.selected[m.items[idx].ID()] {
					allSel = false
					break
				}
			}
			for _, idx := range m.filtered {
				m.selected[m.items[idx].ID()] = !allSel
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m pickerModel) selCount() int {
	n := 0
	for _, v := range m.selected {
		if v {
			n++
		}
	}
	return n
}

func (m pickerModel) View() string {
	if !m.ready {
		return "loading…"
	}
	listW := m.width/2 - 2

	header := titleSty.Render("📓 commit-chronicle") + "  " +
		dimSty.Render(fmt.Sprintf("%s · %s · %d/%d shown · %d selected",
			m.author, m.rangeLabel, len(m.filtered), len(m.items), m.selCount()))

	rows := m.listRows()
	var lines []string
	end := min(m.top+rows, len(m.filtered))
	for i := m.top; i < end; i++ {
		it := m.items[m.filtered[i]]
		box := "[ ]"
		if m.selected[it.ID()] {
			box = "[x]"
		}
		// Plain text first so truncation counts real characters.
		line := fmt.Sprintf("%s %-7s %s  %-16s %-8s %s",
			box, it.Tag(), it.Date, truncate(it.RepoName, 16), it.Ref(), it.Title)
		line = truncate(line, listW-3)
		switch {
		case i == m.cursor:
			line = curSty.Render("❯ " + line)
		case m.selected[it.ID()]:
			line = selSty.Render("  " + line)
		default:
			line = "  " + line
		}
		lines = append(lines, line)
	}
	for len(lines) < rows {
		lines = append(lines, "")
	}
	listBox := boxSty.Width(listW).Height(rows).Render(strings.Join(lines, "\n"))
	prevBox := boxSty.Width(m.preview.Width).Height(rows).Render(m.preview.View())
	body := lipgloss.JoinHorizontal(lipgloss.Top, listBox, prevBox)

	filterLine := dimSty.Render("filter: ") + m.filter.View()
	if !m.filtering && m.filter.Value() == "" {
		filterLine = dimSty.Render("press / to filter")
	}
	help := dimSty.Render("↑↓ move · space/tab select · a all · / filter · enter confirm · q cancel")

	return strings.Join([]string{header, filterLine, body, help}, "\n")
}

func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}
