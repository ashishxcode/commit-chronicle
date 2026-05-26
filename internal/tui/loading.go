package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashishxcode/commit-chronicle/internal/model"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WorkFn does the gathering, reporting progress via the report callback.
type WorkFn func(report func(stage string, n int)) ([]model.Item, error)

type statusMsg struct {
	stage string
	n     int
}
type doneMsg struct {
	items []model.Item
	err   error
}

type loadingModel struct {
	sp       spinner.Model
	label    string
	done     []string // completed phases, e.g. "git history: 87"
	active   string   // phase currently running, e.g. "PRs you reviewed"
	total    int
	start    time.Time
	finished bool
	items    []model.Item
	err      error
}

// RunWithSpinner shows an animated spinner with live progress while work runs
// in the background, then returns its result.
func RunWithSpinner(label string, work WorkFn) ([]model.Item, error) {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(cBlue)

	m := loadingModel{sp: sp, label: label, start: time.Now()}
	p := tea.NewProgram(m)

	go func() {
		items, err := work(func(stage string, n int) {
			p.Send(statusMsg{stage: stage, n: n})
		})
		p.Send(doneMsg{items: items, err: err})
	}()

	res, err := p.Run()
	if err != nil {
		return nil, err
	}
	fm := res.(loadingModel)
	return fm.items, fm.err
}

func (m loadingModel) Init() tea.Cmd { return m.sp.Tick }

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		if msg.n == 0 {
			m.active = msg.stage // phase started
		} else {
			m.done = append(m.done, fmt.Sprintf("%s: %d", msg.stage, msg.n))
			m.total += msg.n
			if m.active == msg.stage {
				m.active = ""
			}
		}
		return m, nil
	case doneMsg:
		m.finished = true
		m.items = msg.items
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.err = fmt.Errorf("canceled")
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.sp, cmd = m.sp.Update(msg)
	return m, cmd
}

func (m loadingModel) View() string {
	if m.finished {
		return "" // cleared once the picker takes over
	}
	elapsed := time.Since(m.start).Truncate(time.Second)
	var b strings.Builder
	fmt.Fprintf(&b, "%s%s  %s\n", m.sp.View(), titleSty.Render(m.label),
		dimSty.Render(fmt.Sprintf("%s · %d found", elapsed, m.total)))

	// Completed phases as a check-list.
	for _, d := range m.done {
		b.WriteString(selSty.Render("  ✓ ") + dimSty.Render(d) + "\n")
	}
	// The phase currently running.
	if m.active != "" {
		b.WriteString(curSty.Render("  → scanning "+m.active) + dimSty.Render(" …"))
	}
	return b.String()
}
