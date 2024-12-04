package cmd

import tea "github.com/charmbracelet/bubbletea"

type SelectedFilePath string
type ContentMsg string
type PathMsg string
type TitleMsg string

func sendSelectedFile(areapath string) tea.Cmd {
	return func() tea.Msg {
		return SelectedFilePath(areapath)
	}
}

func sendContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return ContentMsg(content)
	}
}

func sendPathCmd(path string) tea.Cmd {
	return func() tea.Msg {
		return PathMsg(path)
	}
}

func sendTitleCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return TitleMsg(content)
	}
}
