package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/list"
)

type Vault struct {
	Dir string
}

func NewVaultFromHome() (*Vault, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %w", err)
	}
	v := &Vault{Dir: filepath.Join(homeDir, ".tnotes")}
	if err := v.Ensure(); err != nil {
		return nil, err
	}
	return v, nil
}

func (v *Vault) Ensure() error {
	if err := os.MkdirAll(v.Dir, 0755); err != nil {
		return fmt.Errorf("error creating vault directory: %w", err)
	}
	return nil
}

func (v *Vault) ListNotes() ([]list.Item, error) {
	entries, err := os.ReadDir(v.Dir)
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		desc := "Markdown note"
		if info, err := entry.Info(); err == nil {
			desc = fmt.Sprintf("Modified: %s (%d bytes)", info.ModTime().Format("02 Jan 15:04"), info.Size())
		}
		items = append(items, noteItem{title: entry.Name(), info: desc})
	}
	return items, nil
}

func (v *Vault) OpenOrCreateNote(name string) (*os.File, string, string, error) {
	safeName := filepath.Base(strings.TrimSpace(name))
	if safeName == "" {
		return nil, "", "", fmt.Errorf("note name cannot be empty")
	}
	if !strings.HasSuffix(safeName, ".md") {
		safeName += ".md"
	}

	filePath := filepath.Join(v.Dir, safeName)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, "", "", err
	}

	existingBytes, err := io.ReadAll(f)
	if err != nil {
		_ = f.Close()
		return nil, "", "", err
	}

	content := string(existingBytes)
	if len(content) > 0 {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
	}

	return f, safeName, content, nil
}

func (v *Vault) SaveNote(f *os.File, content string) error {
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return f.Sync()
}

func (v *Vault) RenameNote(oldName, newName string) (string, error) {
	safeOld := filepath.Base(strings.TrimSpace(oldName))
	safeNew := filepath.Base(strings.TrimSpace(newName))
	if safeNew == "" {
		return "", fmt.Errorf("note name cannot be empty")
	}
	if !strings.HasSuffix(safeNew, ".md") {
		safeNew += ".md"
	}

	oldPath := filepath.Join(v.Dir, safeOld)
	newPath := filepath.Join(v.Dir, safeNew)

	if oldPath != newPath {
		if _, err := os.Stat(newPath); err == nil {
			return "", fmt.Errorf("a note named %s already exists", safeNew)
		}
		if err := os.Rename(oldPath, newPath); err != nil {
			return "", err
		}
	}

	return safeNew, nil
}

func (v *Vault) DeleteNote(name string) error {
	safeName := filepath.Base(strings.TrimSpace(name))
	return os.Remove(filepath.Join(v.Dir, safeName))
}

type noteItem struct {
	title, info string
}

func (i noteItem) Title() string       { return i.title }
func (i noteItem) Description() string { return i.info }
func (i noteItem) FilterValue() string { return i.title }
