package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/synaptrix8-cell/tnotes/internal/storage"
	"github.com/synaptrix8-cell/tnotes/internal/tui"
)

func main() {
	vault, err := storage.NewVaultFromHome()
	if err != nil {
		fmt.Printf("There is an error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.NewModel(vault))
	if _, err := p.Run(); err != nil {
		fmt.Printf("There is an error: %v\n", err)
		os.Exit(1)
	}
}
