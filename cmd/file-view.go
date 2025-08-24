package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/The-Mod-Elephant/infinity_dialog/pkg/readers"
	"github.com/The-Mod-Elephant/infinity_file_formats/bg"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Fileview struct {
	title    string
	content  string
	viewport viewport.Model
}

func NewFileView() Fileview {
	f := Fileview{}
	headerHeight := lipgloss.Height(f.headerView())
	footerHeight := lipgloss.Height(f.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	f.viewport = viewport.New(width, height-verticalMarginHeight)
	f.viewport.YPosition = headerHeight
	return f
}

func GetFileContents(path string) (string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	switch strings.ToLower(filepath.Ext(info.Name())) {
	case ".are":
		area, err := bg.OpenArea(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = area.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".bam":
		bam, err := bg.OpenBAM(f, nil)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = bam.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".cre":
		cre, err := bg.OpenCre(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = cre.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".dlg":
		dlg, err := bg.OpenDlg(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = dlg.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".eff":
		effv1, effv2, err := bg.OpenEff(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		if effv1 != nil {
			err = effv1.WriteJson(buf)
		} else {
			err = effv2.WriteJson(buf)
		}
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".itm":
		item, err := bg.OpenITM(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = item.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".sto":
		dlg, err := bg.OpenSTO(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = dlg.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case ".spl":
		dlg, err := bg.OpenSPL(f)
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		err = dlg.WriteJson(buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	default:
		if contents, err := readers.ReadFileToString(path); err != nil {
			return contents, nil
		}
	}
	return "", nil
}

func (f Fileview) Init() tea.Cmd {
	return nil
}

func (f Fileview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TitleMsg:
		f.title = string(msg)
		return f, nil
	case ContentMsg:
		f.content = string(msg)
		f.viewport.SetContent(f.content)
		return f, f.Init()
	case SelectedFilePath:
		content, err := GetFileContents(string(msg))
		if err != nil {
			return f, tea.Quit
		}
		f.content = content
		f.viewport.SetContent(f.content)
		f.title = filepath.Base(string(msg))
		return f, f.Init()
	case PathMsg:
		content, err := GetFileContents(string(msg))
		if err != nil {
			return f, tea.Quit
		}
		f.content = content
		f.viewport.SetContent(f.content)
		f.title = filepath.Base(string(msg))
		return f, f.Init()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return state.PreviousCommand(), nil
		case "ctrl+c", "ctrl+d":
			return f, tea.Quit
		}
	case tea.WindowSizeMsg:
		setViewport(f, msg)
	}
	var cmd tea.Cmd
	f.viewport, cmd = f.viewport.Update(msg)
	return f, cmd
}

func setViewport(f Fileview, msg tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height(f.headerView())
	footerHeight := lipgloss.Height(f.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	f.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
	f.viewport.YPosition = headerHeight
	f.viewport.SetContent(f.content)
	f.viewport.Width = msg.Width
	f.viewport.Height = msg.Height - verticalMarginHeight
}

func (f Fileview) View() string {
	return fmt.Sprintf("%s\n%s\n%s", f.headerView(), f.viewport.View(), f.footerView())
}

func (f Fileview) headerView() string {
	title := titleStyle.Render(f.title)
	line := strings.Repeat("─", max(0, f.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (f Fileview) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", f.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, f.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
