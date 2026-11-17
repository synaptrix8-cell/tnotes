package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/synaptrix8-cell/tnotes/internal/storage"
	appstyles "github.com/synaptrix8-cell/tnotes/internal/styles"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	vault                  *storage.Vault
	newFileInput           textinput.Model
	msg                    string
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
	noteList               list.Model
	noteListVisible        bool
	confirmDelete          bool
	deleteTarget           string
	renaming               bool
	renameTarget           string
	renameInput            textinput.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func NewModel(vault *storage.Vault) tea.Model {
	// Initialize create-note input.
	ti := textinput.New()
	ti.Placeholder = "write the name of your note.."
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(50)
	ti.Prompt = "-->> "
	inStyles := appstyles.NewInputStyles()
	ti.SetStyles(inStyles)

	// Initialize rename input.
	ri := textinput.New()
	ri.Placeholder = "new note name..."
	ri.CharLimit = 156
	ri.SetWidth(50)
	ri.Prompt = "-->> "
	ri.SetStyles(inStyles)

	// Initialize note editor.
	ta := textarea.New()
	ta.Placeholder = "write your note here..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.Prompt = "... "
	ta.SetStyles(appstyles.NewTextAreaStyles())

	// Initialize note list.
	nl := list.New([]list.Item{}, list.NewDefaultDelegate(), 50, 12)
	nl.Title = "Your Notes"
	nl.SetShowHelp(false)

	return model{
		vault:                  vault,
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		noteList:               nl,
		noteListVisible:        false,
		renameInput:            ri,
	}
}

func (m *model) refreshNoteList() error {
	items, err := m.vault.ListNotes()
	if err != nil {
		return err
	}
	m.noteList.SetItems(items)
	return nil
}

func (m *model) closeCurrentFile() {
	if m.currentFile != nil {
		_ = m.currentFile.Close()
		m.currentFile = nil
	}
}

func (m *model) openForEditing(name string) (tea.Model, tea.Cmd) {
	f, safeName, content, err := m.vault.OpenOrCreateNote(name)
	if err != nil {
		m.msg = fmt.Sprintf("Error opening note: %v", err)
		return *m, nil
	}

	m.closeCurrentFile()
	m.currentFile = f
	m.noteListVisible = false
	m.createFileInputVisible = false
	m.newFileInput.Reset()
	m.newFileInput.Blur()

	m.noteTextArea.Reset()
	if len(content) > 0 {
		m.noteTextArea.SetValue(content)
		m.noteTextArea.CursorEnd()
		m.msg = fmt.Sprintf("Opened existing note: %s", safeName)
	} else {
		m.msg = fmt.Sprintf("Created new note: %s", safeName)
	}
	m.noteTextArea.Focus()
	return *m, textarea.Blink
}

func (m *model) selectedNoteTitle() (string, bool) {
	selected := m.noteList.SelectedItem()
	if selected == nil {
		return "", false
	}
	it, ok := selected.(interface{ FilterValue() string })
	if !ok {
		return "", false
	}
	return strings.TrimSpace(it.FilterValue()), true
}
