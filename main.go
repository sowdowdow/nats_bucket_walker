package main

import (
	"fmt"
	natsbinding "nats_bucket_walker/nats_binding"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table    table.Model
	inBucket bool
}

func ListBuckets(model *model) {
	// data retrieval
	buckets, err := natsbinding.GetAllBuckets()
	if err != nil {
		panic(err)
	}
	newRows := []table.Row{}
	for _, b := range buckets {
		newRows = append(newRows, table.Row{b})
	}

	model.table.SetRows(newRows)
	model.inBucket = false
	model.table.Columns()[0].Title = "Buckets"
}

func OpenBucket(model *model) tea.Cmd {
	model.inBucket = true
	selected := model.table.SelectedRow()[0]
	kvs, err := natsbinding.GetAllKV(selected)
	if err != nil {
		if err.Error() == "no keys found" {
			ListBuckets(model)
			return tea.Printf(err.Error())
		}
		panic(err)
	}

	newRows := []table.Row{}
	for _, kv := range kvs {
		newRows = append(newRows, table.Row{kv})
	}

	model.table.SetRows(newRows)
	model.table.Columns()[0].Title = selected
	model.inBucket = true

	return tea.Printf("Opening %s", selected)
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 5)
		m.table.SetWidth(msg.Width - 2)

		m.table.Columns()[0].Width = msg.Width - 4
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "left", "h":
			if m.inBucket {
				ListBuckets(&m) // getting out of the bucket
				return m, tea.Printf("Quitting bucket")
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter", "right", "l":
			if m.inBucket {
				return m, tea.Printf("Already in a bucket")
			}
			cmd := OpenBucket(&m)
			return m, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func main() {

	// data retrieval
	buckets, err := natsbinding.GetAllBuckets()
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "Bucket", Width: 30},
	}

	rows := []table.Row{}
	for _, b := range buckets {
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

	m := model{t, false}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
