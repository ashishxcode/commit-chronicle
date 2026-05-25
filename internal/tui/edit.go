package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type editModel struct {
	ta       textarea.Model
	canceled bool
	width    int
	height   int
}

// Edit opens an in-app editor pre-filled with text.
// Ctrl+S / Ctrl+D saves; Esc cancels. Returns the edited text.
func Edit(text string) (string, bool, error) {
	ta := textarea.New()
	ta.SetValue(text)
	ta.CharLimit = 0
	ta.ShowLineNumbers = true
	// SetValue leaves the cursor at the end; start at the top instead.
	for ta.Line() > 0 {
		ta.CursorUp()
	}
	ta.CursorStart()
	ta.Focus()

	m := editModel{ta: ta}
	res, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return text, false, err
	}
	fm := res.(editModel)
	if fm.canceled {
		return text, true, nil
	}
	return fm.ta.Value(), false, nil
}

func (m editModel) Init() tea.Cmd { return textarea.Blink }

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ta.SetWidth(msg.Width - 2)
		m.ta.SetHeight(msg.Height - 2)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "ctrl+d":
			return m, tea.Quit
		case "esc":
			m.canceled = true
			return m, tea.Quit
		}
	}
	m.ta, cmd = m.ta.Update(msg)
	return m, cmd
}

func (m editModel) View() string {
	help := dimSty.Render("ctrl+s save · esc cancel")
	return strings.Join([]string{m.ta.View(), help}, "\n")
}
