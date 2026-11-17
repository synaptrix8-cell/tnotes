package styles

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

var (
	DocStyle = lipgloss.NewStyle().Margin(1, 2)
)

func WelcomeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2)
}

func WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))
}

func NewInputStyles() textinput.Styles {
	ti := textinput.New()
	st := ti.Styles()
	st.Focused.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	st.Focused.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
	st.Focused.Placeholder = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	st.Cursor.Color = lipgloss.Color("205")
	return st
}

func NewTextAreaStyles() textarea.Styles {
	st := textarea.DefaultStyles(true)
	st.Focused.CursorLine = lipgloss.NewStyle().Background(lipgloss.NoColor{})
	st.Blurred.CursorLine = lipgloss.NewStyle().Background(lipgloss.NoColor{})
	return st
}
