package cmd

import (
	"errors"
	"os"
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
	selectedFile string
	quitting     bool
	err          error
}

func NewDirectoryPicker() directoryPicker {
	fp := filepicker.New()
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.Height = 5
	fp.CurrentDirectory, _ = os.Getwd()
	return directoryPicker{
		filepicker: fp,
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
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			if d.selectedFile != "" {
				l := NewList(d.selectedFile)
				return l, l.Init()
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
		s.WriteString("Pick a directory:\n")
	} else {
		s.WriteString("Selected directory: " + d.filepicker.Styles.Selected.Render(d.selectedFile))
	}

	s.WriteString("\n\n" + d.filepicker.View() + "\n")
	s.WriteString(helpStyle.Render("Keys: e->done, escape->./.., q->quit"))
	return s.String()
}
