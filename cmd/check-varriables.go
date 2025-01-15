package cmd

import (
	"cmp"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strconv"
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

func difference(slice1 *[]string, slice2 *[]string) *[]string {
	diff := []string{}
	m := map[string]int{}

	for _, s := range *slice1 {
		m[s] = 1
	}
	for _, s := range *slice2 {
		m[s] += 1
	}
	for k, v := range m {
		if v == 1 {
			diff = append(diff, k)
		}
	}
	slices.SortFunc(diff, func(a, b string) int {
		v1, _ := strconv.Atoi(a)
		v2, _ := strconv.Atoi(b)
		return cmp.Compare(v1, v2)
	})
	return &diff
}

func (c checkVariables) genRows() *[]table.Row {
	rows := map[string]map[string][]string{}
	_ = filepath.WalkDir(c.root, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && strings.ToLower(ext) == ".tra" {
			lang := strings.ToLower(filepath.Base(filepath.Dir(path)))
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
	for lang, _ := range rows {
		for filename, stringVariables := range largest {
			size_for_lang := rows[lang][filename]
			diff := strings.Join(*difference(&stringVariables, &size_for_lang), ",")
			out = append(out, table.Row{lang, filename, diff})
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
			lang := c.table.SelectedRow()[0]
			file_name := c.table.SelectedRow()[1]
			strings := strings.Split(c.table.SelectedRow()[2], ",")
			path := c.loadFiles[lang][file_name]
			content := []string{"\n"}
			for _, missing := range strings {
				content = append(content, fmt.Sprintf("@%s = ~~\n", missing))
			}
			util.WriteToFile(path, &content)
		case "e", "enter":
			lang := c.table.SelectedRow()[0]
			file_name := c.table.SelectedRow()[1]
			return state.SetAndGetNextCommand(c), SendPathCmd(c.loadFiles[lang][file_name])
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
