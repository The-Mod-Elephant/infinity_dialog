package cmd

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/The-Mod-Elephant/infinity_dialog/pkg/readers"
	"github.com/The-Mod-Elephant/infinity_dialog/pkg/translation"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ListVariables struct {
	table table.Model
}

func generateRows(path string, file fs.FileInfo) *[]table.Row {
	rows := []table.Row{}
	fileContent, err := readers.ReadFileToSlice(path)
	if err != nil {
		return &rows
	}
	variables, err := translation.FromFileContents(fileContent)
	if err == nil {
		lang := filepath.Base(filepath.Dir(path))
		for _, v := range *variables {
			rows = append(rows, table.Row{file.Name(), lang, v.Identifier, v.Value})
		}
	}
	return &rows
}

func NewList() ListVariables {
	columns := []table.Column{
		{Title: "FileName", Width: int(0.2 * float64(width))},
		{Title: "Lang", Width: int(0.1 * float64(width))},
		{Title: "Id", Width: int(0.05 * float64(width))},
		{Title: "Value", Width: int(0.55 * float64(width))},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(height-7),
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

	return ListVariables{table: t}
}

func (l ListVariables) readPath(path string) (*[]table.Row, error) {
	rows := []table.Row{}
	err := filepath.WalkDir(path, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && strings.EqualFold(ext, ".tra") {
			info, _ := file.Info()
			fileRows := *generateRows(path, info)
			rows = append(rows, fileRows...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

func (l ListVariables) Init() tea.Cmd { return nil }

func (l ListVariables) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SelectedFilePath:
		rows, err := l.readPath(string(msg))
		if err != nil {
			return l, tea.Quit
		}
		l.table.SetRows(*rows)
		return l, nil
	case tea.WindowSizeMsg:
		h, w := docStyle.GetFrameSize()
		h1, w1 := baseStyle.GetFrameSize()
		h += h1
		w += w1
		if msg.Height > h {
			l.table.SetHeight(msg.Height - h)
		}
		if msg.Width > w {
			ratio := float64(msg.Width - w)
			l.table.SetColumns([]table.Column{
				{Title: "FileName", Width: int(0.2 * ratio)},
				{Title: "Lang", Width: int(0.1 * ratio)},
				{Title: "Id", Width: int(0.05 * ratio)},
				{Title: "Value", Width: int(0.55 * ratio)},
			})
			l.table.SetWidth(int(ratio))
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return state.PreviousCommand(), nil
		case "ctrl+c", "ctrl+d":
			return l, tea.Quit
		case "enter":
			content := ""
			counter := 0
			for _, s := range strings.Split(l.table.SelectedRow()[3], " ") {
				counter += len(s)
				if counter >= 42 {
					content += "\n"
					counter = 0
				} else if strings.Trim(content, " ") != "" {
					content += " "
				}
				content += s
			}
			title := strings.Join(l.table.SelectedRow()[:3], " ")
			return state.SetAndGetNextCommand(l), tea.Sequence(SendTitleCmd(title), SendContentCmd(content))
		}
	}
	var cmd tea.Cmd
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}

func (l ListVariables) View() string {
	body := []string{l.table.View(), "\n\n", l.table.HelpView(), " enter"}
	return baseStyle.Render(body...)
}
