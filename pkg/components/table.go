package components

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func DynamicalSetTableSize(t *table.Model, msg *tea.WindowSizeMsg, h int, w int) {
	if msg.Height > h {
		t.SetHeight(msg.Height - h)
	}
	if msg.Width > w {
		ratio := float64(msg.Width - w)
		for _, c := range t.Columns() {
			c.Width = int((float64(t.Width()) / float64(c.Width)) * ratio)
		}
		t.SetWidth(int(ratio))
	}
}
