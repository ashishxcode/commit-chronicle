package tui

import (
	"fmt"
	"strings"

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
	sp     spinner.Model
	label  string
	stages []string
	total  int
	done   bool
	items  []model.Item
	err    error
}

// RunWithSpinner shows an animated spinner with live progress while work runs
// in the background, then returns its result.
func RunWithSpinner(label string, work WorkFn) ([]model.Item, error) {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(cBlue)

	m := loadingModel{sp: sp, label: label}
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
		if msg.n > 0 {
			m.stages = append(m.stages, fmt.Sprintf("%s: %d", msg.stage, msg.n))
			m.total += msg.n
		}
		return m, nil
	case doneMsg:
		m.done = true
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
	if m.done {
		return "" // cleared once the picker takes over
	}
	head := m.sp.View() + titleSty.Render(m.label)
	if len(m.stages) == 0 {
		return head + dimSty.Render("  …")
	}
	return head + "\n" + dimSty.Render("  "+strings.Join(m.stages, " · "))
}
