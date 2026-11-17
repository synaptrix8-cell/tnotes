package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	appstyles "github.com/synaptrix8-cell/tnotes/internal/styles"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appstyles.DocStyle.GetFrameSize()
		m.noteList.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
		return m, cmd
	}
	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
		return m, cmd
	}
	if m.noteListVisible {
		m.noteList, cmd = m.noteList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" || msg.String() == "ctrl+q" {
		return m, tea.Quit
	}

	if m.confirmDelete {
		return m.handleConfirmDelete(msg)
	}
	if m.renaming {
		return m.handleRenaming(msg)
	}

	switch msg.String() {
	case "q":
		if !m.createFileInputVisible && m.currentFile == nil && !m.noteListVisible {
			return m, tea.Quit
		}
	case "d":
		if m.noteListVisible && m.noteList.FilterState() != list.Filtering {
			if selected, ok := m.selectedNoteTitle(); ok {
				m.confirmDelete = true
				m.deleteTarget = selected
				m.msg = ""
				return m, nil
			}
		}
	case "r":
		if m.noteListVisible && m.noteList.FilterState() != list.Filtering {
			if selected, ok := m.selectedNoteTitle(); ok {
				m.renaming = true
				m.renameTarget = selected
				m.renameInput.SetValue(strings.TrimSuffix(selected, ".md"))
				m.renameInput.CursorEnd()
				m.renameInput.Focus()
				m.msg = ""
				return m, textinput.Blink
			}
		}
	case "ctrl+n":
		m.msg = "Creating a new note..."
		m.createFileInputVisible = true
		m.newFileInput.Focus()
		return m, textinput.Blink
	case "ctrl+s":
		return m.handleSave()
	case "enter":
		return m.handleEnter()
	case "ctrl+l":
		return m.handleListNotes()
	case "esc":
		if m.noteListVisible && m.noteList.FilterState() == list.Filtering {
			break
		}
		m.msg = "Returned to main menu."
		m.createFileInputVisible = false
		m.noteListVisible = false
		m.newFileInput.Reset()
		m.newFileInput.Blur()
		m.closeCurrentFile()
		return m, nil
	}

	return m, nil
}

func (m model) handleConfirmDelete(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if err := m.vault.DeleteNote(m.deleteTarget); err != nil {
			m.msg = fmt.Sprintf("Error deleting note: %v", err)
		} else {
			m.msg = fmt.Sprintf("Deleted: %s", m.deleteTarget)
		}
		m.confirmDelete = false
		m.deleteTarget = ""

		if err := m.refreshNoteList(); err == nil && len(m.noteList.Items()) == 0 {
			m.noteListVisible = false
			m.msg += " (No notes remaining)"
		}
		return m, nil

	case "n", "N", "esc":
		m.confirmDelete = false
		m.deleteTarget = ""
		m.msg = "Deletion canceled."
		return m, nil
	}
	return m, nil
}

func (m model) handleRenaming(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		newName := strings.TrimSpace(m.renameInput.Value())
		if newName == "" {
			m.msg = "Note name cannot be empty."
			return m, nil
		}

		safeName, err := m.vault.RenameNote(m.renameTarget, newName)
		if err != nil {
			m.msg = fmt.Sprintf("Error renaming note: %v", err)
			return m, nil
		}

		m.msg = fmt.Sprintf("Renamed %s -> %s", m.renameTarget, safeName)
		m.renaming = false
		m.renameTarget = ""
		m.renameInput.Reset()
		m.renameInput.Blur()
		_ = m.refreshNoteList()
		return m, nil

	case "esc":
		m.renaming = false
		m.renameTarget = ""
		m.renameInput.Reset()
		m.renameInput.Blur()
		m.msg = "Rename canceled."
		return m, nil
	}

	var cmd tea.Cmd
	m.renameInput, cmd = m.renameInput.Update(msg)
	return m, cmd
}

func (m model) handleSave() (tea.Model, tea.Cmd) {
	m.msg = "Saving the current note..."
	if m.currentFile != nil {
		if err := m.vault.SaveNote(m.currentFile, m.noteTextArea.Value()); err != nil {
			m.msg = fmt.Sprintf("Error saving note: %v", err)
			return m, nil
		}
		m.msg = fmt.Sprintf("Saved: %s", filepath.Base(m.currentFile.Name()))
	}
	m.closeCurrentFile()
	m.noteTextArea.SetValue("")
	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	if m.createFileInputVisible {
		fileName := strings.TrimSpace(m.newFileInput.Value())
		if fileName == "" {
			m.msg = "Note name cannot be empty."
			return m, nil
		}
		return m.openForEditing(fileName)
	}

	if m.noteListVisible && m.noteList.FilterState() != list.Filtering {
		if selected, ok := m.selectedNoteTitle(); ok {
			return m.openForEditing(selected)
		}
	}

	return m, nil
}

func (m model) handleListNotes() (tea.Model, tea.Cmd) {
	m.msg = "Listing all notes..."
	if err := m.refreshNoteList(); err != nil {
		m.msg = fmt.Sprintf("Error listing notes: %v", err)
		return m, nil
	}
	if len(m.noteList.Items()) == 0 {
		m.msg = "No notes found. Press Ctrl+N to create one."
		return m, nil
	}

	m.closeCurrentFile()
	m.createFileInputVisible = false
	m.newFileInput.Blur()
	m.noteListVisible = true
	m.msg = ""
	return m, nil
}
