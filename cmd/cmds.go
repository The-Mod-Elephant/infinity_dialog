package cmd

import tea "github.com/charmbracelet/bubbletea"

type (
	SelectedFilePath string
	ContentMsg       string
	PathMsg          string
	TitleMsg         string
)

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
