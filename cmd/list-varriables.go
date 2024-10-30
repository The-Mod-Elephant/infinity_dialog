package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dark0dave/infinity_dialog/pkg/translation"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type listVariables struct {
	heading []string
	table   table.Model
}

func generateRows(path string, file fs.FileInfo) *[]table.Row {
	rows := []table.Row{}
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return &rows
	}
	variables, err := translation.FromFileContents(string(fileContent))
	if err == nil {
		for _, v := range *variables {
			rows = append(rows, table.Row{file.Name(), v.Identifier, v.Value})
		}
	}
	return &rows
}

func NewList(path string) listVariables {
	columns := []table.Column{
		{Title: "FileName", Width: 12},
		{Title: "Id", Width: 4},
		{Title: "Value", Width: 40},
	}
	rows := []table.Row{}
	filepath.WalkDir(path, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && ext == ".tra" {
			info, _ := file.Info()
			file_rows := *generateRows(path, info)
			rows = append(rows, file_rows...)
		}

		return nil
	})

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
		BorderTop(true).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return listVariables{table: t}
}

func (l listVariables) Init() tea.Cmd { return nil }

func (l listVariables) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			l.heading = []string{""}
			counter := 0
			for _, s := range strings.Split(l.table.SelectedRow()[2], " ") {
				counter += len(s)
				if counter >= 42 {
					l.heading = append(l.heading, "\n")
					counter = 0
				}
				l.heading = append(l.heading, s)
			}
		}
	}
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}

func (l listVariables) View() string {
	body := []string{l.table.View(), "\n\n", l.table.HelpView(), " enter"}
	if len(l.heading) == 0 {
		return baseStyle.Render(body...)
	} else {
		payload := l.heading
		payload = append(payload, "\n")
		payload = append(payload, body...)
		return baseStyle.Render(payload...)
	}
}
