package main

import (
	natsbinding "nats_bucket_walker/nats"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) ListBuckets() tea.Cmd {
	// data retrieval
	buckets, err := natsbinding.GetAllBuckets()
	if err != nil {
		panic(err)
	}
	newRows := []table.Row{}
	for _, b := range buckets {
		newRows = append(newRows, table.Row{b})
	}

	m.table.SetRows(newRows)
	m.inBucket = false
	m.table.Columns()[0].Title = "Bucket"
	return tea.Batch(
		tea.Printf("Quitting bucket"),
	)
}

func (m *model) OpenBucket() tea.Cmd {

	m.inBucket = true
	name := m.table.SelectedRow()[0]
	m.bucket = bucket{watching: false, name: name}
	kvs, err := natsbinding.GetAllKV(name)
	if err != nil {
		panic(err)
	}

	newRows := []table.Row{}
	for _, kv := range kvs {
		newRows = append(newRows, table.Row{kv})
	}

	m.table.SetRows(newRows)
	m.table.Columns()[0].Title = name
	m.inBucket = true

	return tea.Batch(
		tea.Printf("Opening %s", name),
	)
}
