package cmd

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dark0dave/infinity_dialog/pkg/translation"
	"github.com/dark0dave/infinity_dialog/pkg/util"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type list struct {
	table table.Model
}

func NewList(path string) list {
	columns := []table.Column{
		{Title: "FileName", Width: 15},
		{Title: "Id", Width: 2},
		{Title: "Value", Width: 35},
	}
	rows := []table.Row{}
	for _, f := range util.GetFiles(path, ".tra") {
		fileContent, err := util.ReadFile(path, f)
		if err != nil {
			continue
		}
		varriables, err := translation.FromFileContents(fileContent)
		if err == nil {
			for _, v := range *varriables {
				rows = append(rows, table.Row{f.Name(), v.Identifier, v.Value})
			}
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return list{table: t}
}

func (l list) Init() tea.Cmd { return nil }

func (l list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if l.table.Focused() {
				l.table.Blur()
			} else {
				l.table.Focus()
			}
		case "q", "ctrl+c":
			return l, tea.Quit
		case "enter":
			return l, tea.Batch(
				tea.Printf("Let's go to %s!", l.table.SelectedRow()[1]),
			)
		}
	}
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}

func (l list) View() string {
	return baseStyle.Render(l.table.View()) + "\n  " + l.table.HelpView() + "\n"
}
