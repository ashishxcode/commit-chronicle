package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chooseModel struct {
	title    string
	options  []string
	cursor   int
	chosen   int
	canceled bool
}

// Choose shows a single-select vertical menu and returns the chosen index.
func Choose(title string, options []string) (int, bool, error) {
	m := chooseModel{title: title, options: options, chosen: -1}
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		return 0, false, err
	}
	fm := res.(chooseModel)
	return fm.chosen, fm.canceled, nil
}

func (m chooseModel) Init() tea.Cmd { return nil }

func (m chooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "ctrl+c", "q", "esc":
			m.canceled = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.chosen = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m chooseModel) View() string {
	var b strings.Builder
	b.WriteString(titleSty.Render(m.title) + "\n\n")
	for i, opt := range m.options {
		if i == m.cursor {
			b.WriteString(curSty.Render("❯ "+opt) + "\n")
		} else {
			b.WriteString(lipgloss.NewStyle().Render("  "+opt) + "\n")
		}
	}
	b.WriteString("\n" + dimSty.Render("↑↓ move · enter select · q cancel"))
	return b.String()
}
