package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected string
}

func InitialModel() model {
	return model{
		choices: []string{
			"Find all strings in a file",
			"Add strings to tra",
			"What range of numbers are free",
		},
		selected: "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k", "w":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j", "s":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.choices[m.cursor]
			d := NewDirectoryPicker()
			return d, d.Init()
		}
	}
	return m, nil
}

func (m model) View() string {
	layout := "What should the command do?\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		checked := " " // not selected
		if len(m.selected) > 0 {
			checked = "x" // selected!
		}
		layout += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	layout += "\nPress q to quit.\n"
	return layout
}
