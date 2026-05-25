package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type customRangeModel struct {
	inputs   []textinput.Model
	focus    int
	canceled bool
	submit   bool
}

// CustomRange prompts for a from/to date range. Returns (from, to, canceled).
// `to` may be empty (caller defaults it to today).
func CustomRange() (string, string, bool, error) {
	from := textinput.New()
	from.Prompt = "From (YYYY-MM-DD): "
	from.Placeholder = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	from.CharLimit = 10
	from.Focus()

	to := textinput.New()
	to.Prompt = "To   (YYYY-MM-DD): "
	to.Placeholder = time.Now().Format("2006-01-02") + " (blank = today)"
	to.CharLimit = 10

	m := customRangeModel{inputs: []textinput.Model{from, to}}
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", "", false, err
	}
	fm := res.(customRangeModel)
	if fm.canceled || !fm.submit {
		return "", "", true, nil
	}
	return strings.TrimSpace(fm.inputs[0].Value()), strings.TrimSpace(fm.inputs[1].Value()), false, nil
}

func (m customRangeModel) Init() tea.Cmd { return textinput.Blink }

func (m customRangeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "ctrl+c", "esc":
			m.canceled = true
			return m, tea.Quit
		case "enter":
			// Submit from the last field; otherwise advance.
			if m.focus == len(m.inputs)-1 {
				m.submit = true
				return m, tea.Quit
			}
			m.focus++
		case "tab", "down":
			m.focus = (m.focus + 1) % len(m.inputs)
		case "shift+tab", "up":
			m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
		}
		for i := range m.inputs {
			if i == m.focus {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}
	}
	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m customRangeModel) View() string {
	var b strings.Builder
	b.WriteString(titleSty.Render("📅 Custom date range") + "\n\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View() + "\n")
	}
	b.WriteString("\n" + dimSty.Render("tab switch · enter next/confirm · esc cancel"))
	return b.String()
}
