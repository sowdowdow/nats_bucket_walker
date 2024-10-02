package main

import (
	"fmt"
	"nats_bucket_walker/core"
	natsbinding "nats_bucket_walker/nats"
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
	bucket   bucket
}

type bucket struct {
	watching bool
	name     string
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
			// TODO : check if bucket is in watch mode

			if m.inBucket {
				m.bucket = bucket{}
				return m, m.ListBuckets()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter", "right", "l":
			if m.inBucket {
				return m, func() tea.Msg {
					// TODO : Add watching
					return msg
				}
			}
			return m, m.OpenBucket()
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

	t := core.InitTable(buckets)
	m := model{t, false, bucket{}}
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
