package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/The-Mod-Elephant/infinity_dialog/pkg/util"
	"github.com/The-Mod-Elephant/infinity_file_formats/bg"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tree "github.com/savannahostrowski/tree-bubble"
)

type nested struct {
	tree      tree.Model
	file_map  *map[string]string
	paginator paginator.Model
}

func (m nested) Init() tea.Cmd {
	return nil
}

func NewTree() nested {
	h, w := docStyle.GetFrameSize()
	_, right, _, left := docStyle.GetPadding()
	w = w - left - right
	h = height - h

	// Paginate
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = height - 5
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	return nested{tree: tree.New([]tree.Node{}, w, h), paginator: p}
}

func (n nested) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PathMsg:
		file_map := map[string]string{}
		filepath.Walk(string(msg), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Only search baf + are files
			if !info.IsDir() && (filepath.Ext(info.Name()) == ".are" || filepath.Ext(info.Name()) == ".baf") {
				file_map[strings.ToLower(info.Name())] = path
			}
			return nil
		})
		n.file_map = &file_map
	case SelectedFilePath:
		nodes := []tree.Node{}
		parseArea(&nodes, string(msg), n.file_map)
		n.paginator.SetTotalPages(size(&nodes))
		n.tree.SetNodes(nodes)
		return n, n.Init()
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		n.tree.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			n.tree.SetCursor((n.paginator.Page + 1) * n.paginator.PerPage)
		case "left":
			if n.tree.Cursor() != 0 {
				n.tree.SetCursor((n.paginator.Page - 1) * n.paginator.PerPage)
			}
		case "up":
			cusor := n.tree.Cursor()
			if cusor != 0 && n.tree.Cursor()%(n.paginator.PerPage) == 0 {
				n.paginator.PrevPage()
			}
		case "down":
			cusor := n.tree.Cursor()
			if cusor != n.paginator.PerPage*n.paginator.TotalPages && (cusor+1)%(n.paginator.PerPage) == 0 {
				n.paginator.NextPage()
			}
		case "e", "enter":
			nodes := n.tree.Nodes()
			node, _ := getSelected(&n, &nodes, 0)
			return state.SetAndGetNextCommand(n), SendPathCmd(node.Desc)
		case "q", "esc":
			return state.SetAndGetPreviousCommand(n), nil
		case "ctrl+c", "ctrl+d":
			return n, tea.Quit
		}
	}
	var (
		pag_cmd  tea.Cmd
		tree_cmd tea.Cmd
	)
	n.tree, tree_cmd = n.tree.Update(msg)
	n.paginator, pag_cmd = n.paginator.Update(msg)
	return n, tea.Batch(pag_cmd, tree_cmd)
}

func (n nested) View() string {
	// TODO: Collapse and expand tree
	items := strings.Split(n.tree.View(), "\n")
	var b strings.Builder
	b.WriteString("\n  Dialogue Tree\n\n")
	start, end := n.paginator.GetSliceBounds(len(items))
	for _, item := range items[start:end] {
		b.WriteString(item + "\n")
	}
	b.WriteString("  " + n.paginator.View())
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return n.tree.Styles.Shapes.Render(b.String())
}

func getSelected(n *nested, nodes *[]tree.Node, counter int) (*tree.Node, int) {
	for _, node := range *nodes {
		counter += 1
		if counter-1 == n.tree.Cursor() {
			return &node, counter
		}
		if len(node.Children) > 0 {
			if child, cnt := getSelected(n, &node.Children, counter); child != nil {
				return child, cnt
			} else {
				counter = cnt
			}
		}
	}
	return nil, counter
}

func parseArea(nodes *[]tree.Node, areapath string, file_map *map[string]string) {
	f, err := os.Open(areapath)
	if err != nil {
		return
	}
	defer f.Close()
	area, err := bg.OpenArea(f)
	if err != nil {
		return
	}

	child_name := fmt.Sprintf("%s.%s", strings.Split(strings.ToLower(string(area.Offsets.Script.Name[:])), "\x00")[0], "baf")
	file_path := (*file_map)[child_name]

	parent := tree.Node{
		Value: filepath.Base(areapath),
		Desc:  areapath,
		Children: []tree.Node{{
			Value:    child_name,
			Desc:     file_path,
			Children: []tree.Node{},
		}},
	}

	*nodes = append((*nodes), parent)
	if err := findChildren(file_path, file_map, nodes, &parent.Children[len(parent.Children)-1], 0); err != nil {
		return
	}

	for _, entrance := range area.Entrances {
		area_name := fmt.Sprintf("%s.%s", strings.ToLower(string(entrance.Name.Value[:])), "are")
		area_path := (*file_map)[child_name]
		if !presentInTopOfTree(*nodes, area_name) {
			parseArea(nodes, area_path, file_map)
		}
	}
}

func findChildren(path string, file_map *map[string]string, nodes *[]tree.Node, child *tree.Node, depth int) error {
	if depth > 3 {
		return nil
	}
	contents, err := util.ReadFileToString(path)
	contents = strings.ToLower(contents)
	if err != nil {
		return err
	}
	filename := filepath.Base(path)
	for k, v := range *file_map {
		if k != filename && strings.Contains(contents, "\""+k[:len(k)-4]+"\")") {
			child.Children = append(child.Children, tree.Node{
				Value:    k,
				Desc:     v,
				Children: []tree.Node{},
			})
			if k[len(k)-3:] == "are" {
				if !presentInTopOfTree(*nodes, k) {
					parseArea(nodes, v, file_map)
				}
			} else {
				if !presentInTreeExcludingTop(nodes, k) {
					err := findChildren(v, file_map, nodes, &child.Children[len(child.Children)-1], depth+1)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func presentInTreeExcludingTop(nodes *[]tree.Node, name string) bool {
	for _, parent := range *nodes {
		for _, child := range parent.Children {
			for _, grandchild := range child.Children {
				if presentInTree(&grandchild.Children, name) {
					return true
				}
			}
		}
	}
	return false
}

func presentInTopOfTree(nodes []tree.Node, name string) bool {
	for _, child := range nodes {
		if child.Value == name {
			return true
		}
	}
	return false
}

func presentInTree(nodes *[]tree.Node, name string) bool {
	for _, child := range *nodes {
		if child.Value == name {
			return true
		}
		if presentInTree(&child.Children, name) {
			return true
		}
	}
	return false
}

func size(nodes *[]tree.Node) int {
	start := 0
	for _, child := range *nodes {
		start += 1
		if len(child.Children) > 0 {
			start += size(&child.Children)
		}
	}
	return start
}
