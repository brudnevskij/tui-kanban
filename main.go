package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	list list.Model
	err  error
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(w, h int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), w, h)
	m.list.Title = "To Do"
	m.list.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawbery milk"},
		Task{status: todo, title: "eat sushi", description: "miso soup"},
		Task{status: todo, title: "cleaning", description: "do laundry"},
	})
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.list.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
