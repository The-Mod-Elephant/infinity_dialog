package cmd

import (
	"os"

	"github.com/The-Mod-Elephant/infinity_dialog/pkg/nav"
	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Title = "Infinity Dialog"
)

var (
	state    = nav.NewState()
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	width    = 0
	height   = 0
)

type Item struct {
	title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type Initial struct {
	list list.Model
}

func InitialModel() Initial {
	items := []list.Item{
		Item{title: "Missing", desc: "Find missing strings for langs in a mod/directory"},
		Item{title: "Discover", desc: "Find all strings in a mod/directory"},
		Item{title: "Traverse", desc: "Show tree of locations through a mod"},
		Item{title: "View", desc: "View any Infinity Engine file or text file"},
		Item{title: "Mods", desc: "View what mods are installed, in a game directory"},
		// TODO: Implement these
		// item{title: "Add", desc: "Add strings to tra"},
		// item{title: "Range", desc: "What range of numbers are free"},
		// item{title: "Convert", desc: "Convert files to be traified"},
		// item{title: "Decompiler", desc: "Dialog decompiler"},
	}
	i := Initial{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	i.list.Title = Title
	return i
}

func (i Initial) Init() tea.Cmd {
	return nil
}

func (i Initial) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			currentPath, err := os.Getwd()
			if err != nil {
				return i, tea.Quit
			}
			switch i.list.SelectedItem().FilterValue() {
			case "Missing":
				d := NewDirectoryPicker(true, "Select a Mod Directory")
				c := NewCheck()
				f := NewFileView()
				state.SetNextCommand(d).SetNextCommand(c).SetNextCommand(f)
				return state.SetAndGetNextCommand(i), SendSelectedFile(currentPath)
			case "Discover":
				d := NewDirectoryPicker(true, "Select a Mod Directory")
				l := NewList()
				f := NewFileView()
				state.SetNextCommand(d).SetNextCommand(l).SetNextCommand(f)
				return state.SetAndGetNextCommand(i), SendSelectedFile(currentPath)
			case "Traverse":
				d := NewDirectoryPicker(true, "Select a Mod Directory")
				f := NewDirectoryPicker(false, "Select an area to start")
				t := NewTree()
				v := NewFileView()
				state.SetNextCommand(d).SetNextCommand(f).SetNextCommand(t).SetNextCommand(v)
				return state.SetAndGetNextCommand(i), SendSelectedFile(currentPath)
			case "View":
				d := NewDirectoryPicker(false, "Select a file to start")
				v := NewFileView()
				state.SetNextCommand(d).SetNextCommand(v)
				return state.SetAndGetNextCommand(i), SendSelectedFile(currentPath)
			case "Mods":
				d := NewDirectoryPicker(true, "Select a game directory folder (BGEE, BG2EE, or EET)")
				m := NewModList()
				state.SetNextCommand(d).SetNextCommand(m)
				return state.SetAndGetNextCommand(i), SendSelectedFile(currentPath)
			}
		}
	}
	var cmd tea.Cmd
	i.list, cmd = i.list.Update(msg)
	return i, cmd
}

func (i Initial) View() string {
	return docStyle.Render(i.list.View())
}
