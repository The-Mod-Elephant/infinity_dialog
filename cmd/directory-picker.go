package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type directoryPicker struct {
	filepicker   filepicker.Model
	nextCommand  func(string) (tea.Model, tea.Cmd)
	message      string
	selectedFile string
	quitting     bool
	err          error
}

func NewDirectoryPicker(dir bool, message string, nextCommand func(string) (tea.Model, tea.Cmd)) directoryPicker {
	fp := filepicker.New()
	fp.DirAllowed = dir
	fp.FileAllowed = !dir
	h, _ := docStyle.GetFrameSize()
	fp.Height = height - h - 5
	fp.AutoHeight = true
	fp.CurrentDirectory = currentDirectory
	return directoryPicker{
		filepicker:  fp,
		nextCommand: nextCommand,
		message:     message,
	}
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (d directoryPicker) Init() tea.Cmd {
	return d.filepicker.Init()
}

func (d directoryPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, _ := docStyle.GetFrameSize()
		d.filepicker.Height = d.filepicker.Height - h
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			if d.selectedFile != "" {
				return d.nextCommand(d.selectedFile)
			}
		case "ctrl+c", "q":
			d.quitting = true
			return d, tea.Quit
		}
	case clearErrorMsg:
		d.err = nil
	}

	var cmd tea.Cmd
	d.filepicker, cmd = d.filepicker.Update(msg)

	if didSelect, path := d.filepicker.DidSelectFile(msg); didSelect {
		d.selectedFile = path
	}

	if didSelect, path := d.filepicker.DidSelectDisabledFile(msg); didSelect {
		d.err = errors.New(path + " is not valid.")
		d.selectedFile = ""
		return d, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return d, cmd
}

func (d directoryPicker) View() string {
	var s strings.Builder
	s.WriteString("\n  ")

	if d.err != nil {
		s.WriteString(d.filepicker.Styles.DisabledFile.Render(d.err.Error()))
	} else if d.selectedFile == "" {
		s.WriteString(d.message + "\n")
	} else {
		fileInfo, _ := os.Stat(d.selectedFile)
		if fileInfo.IsDir() {
			currentDirectory = d.selectedFile
		} else {
			currentDirectory = filepath.Dir(d.selectedFile)
		}
		s.WriteString("Selected: " + d.filepicker.Styles.Selected.Render(d.selectedFile))
	}

	s.WriteString("\n\n" + d.filepicker.View() + "\n")
	s.WriteString(helpStyle.Render("  Keys: e->done, escape->./.., q->quit"))
	return s.String()
}
