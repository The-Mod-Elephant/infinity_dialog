package cmd

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dark0dave/infinity_dialog/pkg/components"
	"github.com/dark0dave/infinity_dialog/pkg/util"
)

type dialogueTree struct {
	table table.Model
}

func NewDialogueTree() dialogueTree {
	columns := []table.Column{
		{Title: "FileName", Width: int(0.2 * float64(width))},
		{Title: "Character Name", Width: int(0.2 * float64(width))},
		{Title: "Value", Width: int(0.55 * float64(width))},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(height-len(columns)+4),
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

	return dialogueTree{table: t}
}

func (d dialogueTree) Init() tea.Cmd {
	return nil
}

func (d dialogueTree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SelectedFilePath:
		// TODO: This only works on banter files
		rows := []table.Row{}
		files, err := util.ReadFiles(string(msg), "d")
		for fileName, fileContents := range files {
			for _, line := range *fileContents {
				if strings.Contains(line, "/*") && (strings.Contains(line, "@") || strings.Contains(line, "==")) {
					var characterName String
					if strings.Contains(line, "@") && strings.Contains(line, "==") {
						characterName = strings.Split(line, " ")[1]
					}

					rows = append(rows, table.Row{
						fileName,
					})
				}
			}
		}
		if err != nil {
			return d, nil
		}
		d.table.SetRows(rows)
	case tea.WindowSizeMsg:
		h, w := docStyle.GetFrameSize()
		components.DynamicalSetTableSize(&d.table, &msg, h, w)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return state.PreviousCommand(), nil
		case "ctrl+c", "ctrl+d":
			return d, tea.Quit
		}
	}
	var cmd tea.Cmd
	d.table, cmd = d.table.Update(msg)
	return d, cmd
}

func (d dialogueTree) View() string {
	return d.table.View()
}
