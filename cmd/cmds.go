package cmd

import tea "github.com/charmbracelet/bubbletea"

type SelectedFilePath string
type ContentMsg string
type PathMsg string
type TitleMsg string

func SendSelectedFile(areapath string) tea.Cmd {
	return func() tea.Msg {
		return SelectedFilePath(areapath)
	}
}

func SendContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return ContentMsg(content)
	}
}

func SendPathCmd(path string) tea.Cmd {
	return func() tea.Msg {
		return PathMsg(path)
	}
}

func SendTitleCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return TitleMsg(content)
	}
}
