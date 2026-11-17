package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	appstyles "github.com/synaptrix8-cell/tnotes/internal/styles"
)

func (m model) View() tea.View {
	welcome := appstyles.WelcomeStyle().Render("welcome to Tnotes")
	warnStyle := appstyles.WarningStyle()

	help := m.helpBar()
	var content string

	if m.msg != "" {
		content += fmt.Sprintf("%s\n\n", m.msg)
	}
	if m.confirmDelete {
		content += warnStyle.Render(fmt.Sprintf("Delete %q? (y/n)", m.deleteTarget)) + "\n\n"
	}
	if m.renaming {
		content += fmt.Sprintf("Rename %q to:\n%s\n\n", m.renameTarget, m.renameInput.View())
	}
	if m.createFileInputVisible {
		content += m.newFileInput.View() + "\n\n"
	}
	if m.noteListVisible {
		content += appstyles.DocStyle.Render(m.noteList.View()) + "\n\n"
	}
	if m.currentFile != nil {
		content += m.noteTextArea.View() + "\n\n"
	}

	v := tea.NewView(fmt.Sprintf("\n%s\n\n%s%s", welcome, content, help))
	v.AltScreen = true
	return v
}

func (m model) helpBar() string {
	if m.confirmDelete {
		return "y: confirm delete | n/esc: cancel | ctrl+q: quit\n"
	}
	if m.renaming {
		return "enter: apply rename | esc: cancel | ctrl+q: quit\n"
	}
	if m.currentFile != nil {
		return "ctrl+s: save | esc: back | ctrl+q: quit\n"
	}
	if m.noteListVisible {
		return "enter: open | d: delete | r: rename | esc: back | ctrl+q: quit\n"
	}
	if m.createFileInputVisible {
		return "enter: create/open | esc: back | ctrl+q: quit\n"
	}
	return "Ctrl+N: New Note\nCtrl+S: Save Note\nCtrl+l: List Notes\nCtrl+Q: Quit\nEsc: back/Save Note\n"
}
