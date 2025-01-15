package cmd

import (
	"os"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dark0dave/infinity_dialog/pkg/nav"
)

var (
	state    = nav.NewState()
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	width    = 0
	height   = 0
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
		item{title: "Check", desc: "Check all strings in a mod/directory"},
		item{title: "Discover", desc: "Find all strings in a mod/directory"},
		item{title: "Traverse", desc: "Show tree of locations through a mod"},
		item{title: "View", desc: "View any Infinity Engine file or text file"},
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
		case "ctrl+c", "ctrl+d", "q":
			return i, tea.Quit
		case "enter", " ":
			state = nav.NewState()
			current_path, err := os.Getwd()
			if err != nil {
				return i, tea.Quit
			}
			switch i.list.SelectedItem().FilterValue() {
			case "Check":
				d := NewDirectoryPicker(true, "Select a Mod Directory")
				c := NewCheck()
				f := NewFileView()
				state.SetNextCommand(d).SetNextCommand(c).SetNextCommand(f)
				return state.SetAndGetNextCommand(i), SendSelectedFile(current_path)
			case "Discover":
				d := NewDirectoryPicker(true, "Select a Mod Directory")
				l := NewList()
				f := NewFileView()
				state.SetNextCommand(d).SetNextCommand(l).SetNextCommand(f)
				return state.SetAndGetNextCommand(i), SendSelectedFile(current_path)
			case "Traverse":
				d := NewDirectoryPicker(true, "Select a Mod Directory")
				f := NewDirectoryPicker(false, "Select an area to start")
				t := NewTree()
				v := NewFileView()
				state.SetNextCommand(d).SetNextCommand(f).SetNextCommand(t).SetNextCommand(v)
				return state.SetAndGetNextCommand(i), SendSelectedFile(current_path)
			case "View":
				d := NewDirectoryPicker(false, "Select a file to start")
				v := NewFileView()
				state.SetNextCommand(d).SetNextCommand(v)
				return state.SetAndGetNextCommand(i), SendSelectedFile(current_path)
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
