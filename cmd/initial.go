package cmd

import (
	"os"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
	currentDirectory, _ = os.Getwd()
	width               = 0
	height              = 0
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type initial struct {
	list list.Model
}

func InitialModel() initial {
	items := []list.Item{
		item{title: "Traverse", desc: "Show tree of locations through a mod"},
		item{title: "Discover", desc: "Find all strings in a mod/directory"},
		// TODO: Implement these
		// item{title: "Add", desc: "Add strings to tra"},
		// item{title: "Range", desc: "What range of numbers are free"},
		// item{title: "Convert", desc: "Convert files to be traified"},
		// item{title: "Decompiler", desc: "Dialogue decompiler"},
	}
	i := initial{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	i.list.Title = "Infinity Dialogue"
	return i
}

func (i initial) Init() tea.Cmd {
	return nil
}

func (i initial) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		i.list.SetSize(msg.Width-h, msg.Height-v)
		height, width = max(msg.Height, height), max(msg.Width, width)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return i, tea.Quit
		case "enter", " ":
			switch i.list.SelectedItem().FilterValue() {
			case "Traverse":
				d := NewDirectoryPicker(true, "Select a Mod Directory", func(dirpath string) (tea.Model, tea.Cmd) {
					f := NewDirectoryPicker(false, "Select an area to start", func(areapath string) (tea.Model, tea.Cmd) {
						t := NewTree(dirpath, areapath)
						return t, t.Init()
					})
					return f, f.Init()
				})
				return d, d.Init()
			case "Discover":
				d := NewDirectoryPicker(true, "Select a Mod Directory", func(path string) (tea.Model, tea.Cmd) {
					l := NewList(path)
					return l, l.Init()
				})
				return d, d.Init()
			}
		}
	}
	var cmd tea.Cmd
	i.list, cmd = i.list.Update(msg)
	return i, cmd
}

func (i initial) View() string {
	return docStyle.Render(i.list.View())
}
