package core

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

func InitTable(data []string) table.Model {

	w, _, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "Bucket", Width: w},
	}

	rows := []table.Row{}
	for _, b := range data {
		rows = append(rows, table.Row{b})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		// table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
