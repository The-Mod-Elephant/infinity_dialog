package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dark0dave/infinity_dialog/pkg/translation"
	"github.com/dark0dave/infinity_dialog/pkg/util"
)

type checkVariables struct {
	table     table.Model
	loadFiles map[string]map[string]string
	root      string
	langDir   string
}

func NewCheck() checkVariables {
	columns := []table.Column{
		{Title: "Lang", Width: int(0.1 * float64(width))},
		{Title: "Filename", Width: int(0.2 * float64(width))},
		{Title: "Missing Ids", Width: int(0.5 * float64(width))},
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

	return checkVariables{table: t}
}

func (c *checkVariables) findPath() string {
	lang := c.table.SelectedRow()[0]
	file_name := c.table.SelectedRow()[1]
	path := c.loadFiles[lang][file_name]
	if len(path) == 0 {
		path = filepath.Join(c.langDir, lang, file_name)
	}
	return path
}

func (c *checkVariables) genRows() *[]table.Row {
	rows := map[string]map[string][]string{}
	_ = filepath.WalkDir(c.root, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && strings.ToLower(ext) == ".tra" {
			if len(c.langDir) == 0 {
				c.langDir = filepath.Dir(filepath.Dir(path))
			}
			lang := filepath.Base(filepath.Dir(path))
			if len(c.loadFiles[lang]) == 0 {
				c.loadFiles[lang] = map[string]string{}
			}
			c.loadFiles[lang][file.Name()] = path
			fileContent, err := util.ReadFileToSlice(path)
			if err != nil {
				return err
			}
			variables, err := translation.FromFileContents(fileContent)
			if err == nil {
				if len(rows[lang]) == 0 {
					rows[lang] = map[string][]string{}
				}
				for _, v := range *variables {
					rows[lang][file.Name()] = append(rows[lang][file.Name()], v.Identifier)
				}
			}
		}
		return nil
	})
	largest := map[string][]string{}
	for _, files := range rows {
		for filename, stringVariables := range files {
			if len(largest[filename]) < len(stringVariables) {
				largest[filename] = stringVariables
			}
		}
	}
	out := []table.Row{}
	for lang := range rows {
		for filename, stringVariables := range largest {
			size_for_lang := rows[lang][filename]
			sliceDiff := util.SortedDifference(&stringVariables, &size_for_lang)
			diff := strings.Join(*sliceDiff, ",")
			if len(diff) > 0 {
				out = append(out, table.Row{lang, filename, diff})
			}
		}
	}
	return &out
}

func (c checkVariables) Init() tea.Cmd { return nil }

func (c checkVariables) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SelectedFilePath:
		c.loadFiles = map[string]map[string]string{}
		c.root = string(msg)
		c.table.SetRows(*c.genRows())
		return c, nil
	case tea.WindowSizeMsg:
		h, w := docStyle.GetFrameSize()
		h1, w1 := baseStyle.GetFrameSize()
		h += h1
		w += w1
		if msg.Height > h {
			c.table.SetHeight(msg.Height - h)
		}
		if msg.Width > w {
			ratio := float64(msg.Width - w)
			c.table.SetColumns([]table.Column{
				{Title: "Lang", Width: int(0.1 * ratio)},
				{Title: "Filename", Width: int(0.2 * float64(width))},
				{Title: "Missing Ids", Width: int(0.5 * float64(width))},
			})
			c.table.SetWidth(int(ratio))
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return state.PreviousCommand(), nil
		case "ctrl+c", "ctrl+d":
			return c, tea.Quit
		case "f":
			if len(c.table.Rows()) > 0 {

				strings := strings.Split(c.table.SelectedRow()[2], ",")
				content := []string{"\n"}
				for _, missing := range strings {
					content = append(content, fmt.Sprintf("@%s = ~~\n", missing))
				}
				err := util.WriteToFile(c.findPath(), &content)
				if err != nil {
					panic(err)
				}
				c.table.SetRows(*c.genRows())
			}
		case "e", "enter":
			if len(c.table.Rows()) > 0 {
				lang := c.table.SelectedRow()[0]
				file_name := c.table.SelectedRow()[1]
				return state.SetAndGetNextCommand(c), SendPathCmd(c.loadFiles[lang][file_name])
			}
		}
	}
	var cmd tea.Cmd
	c.table, cmd = c.table.Update(msg)
	return c, cmd
}

func (c checkVariables) View() string {
	body := []string{c.table.View(), "\n\n", c.table.HelpView(), " e enter view, f fix"}
	return baseStyle.Render(body...)
}
