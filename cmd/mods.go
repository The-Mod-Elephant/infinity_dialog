package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/The-Mod-Elephant/infinity_file_formats/bg"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ModList struct {
	*table.Model
}

func NewModList() ModList {
	columns := []table.Column{
		{Title: "Index", Width: int(0.1 * float64(width))},
		{Title: "Tp2 file", Width: int(0.35 * float64(width))},
		{Title: "Language", Width: int(0.1 * float64(width))},
		{Title: "Component", Width: int(0.1 * float64(width))},
		{Title: "Component name", Width: int(0.35 * float64(width))},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(height-7),
	)
	return ModList{&t}
}

func (m ModList) Init() tea.Cmd {
	return nil
}

func (m ModList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.Focused() {
		return m, nil
	}

	switch msg := msg.(type) {
	case SelectedFilePath:
		path := filepath.Join(filepath.Clean(string(msg)), "weidu.log")
		file, err := os.Open(filepath.Clean(path))
		if err != nil {
			return m, tea.Quit
		}
		defer file.Close()
		log, err := bg.OpenLog(file)
		if err != nil {
			return m, tea.Quit
		}
		rows := []table.Row{}
		for i, c := range log.Components {
			items := []string{strconv.Itoa(i), fmt.Sprintf("%s%s%s", c.Name, string(os.PathSeparator), c.TpFile), c.Lang, c.Component, c.ComponentName}
			rows = append(rows, table.Row(items))
		}
		m.SetRows(rows)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.Height())
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.Height())
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.Height() / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.Height() / 2)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		}
		switch msg.String() {
		case "q", "esc":
			return state.PreviousCommand(), nil
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}

	return m, nil
}
