package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

/* MODEL MANAGEMENT*/
var models []tea.Model

const (
	modelMain status = iota
	modelForm
)

/* STYLING */
var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(lipgloss.Color("241"))
)

type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) *Task {
	return &Task{
		status:      status,
		title:       title,
		description: description,
	}
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
		return
	}
	t.status++
}

func (t *Task) FilterValue() string {
	return t.title
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveItemToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(*Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items()), selectedItem)
	return nil
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
		return
	}
	m.focused++
}

func (m *Model) Previous() {
	if m.focused == todo {
		m.focused = done
		return
	}
	m.focused--
}

func (m *Model) initLists(w, h int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), w/divisor, h-divisor/2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	// Init To DO
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		&Task{status: todo, title: "buy milk", description: "strawbery milk"},
		&Task{status: todo, title: "eat sushi", description: "miso soup"},
		&Task{status: todo, title: "cleaning", description: "do laundry"},
	})

	m.lists[inProgress].Title = "In progress"
	m.lists[inProgress].SetItems([]list.Item{
		&Task{status: inProgress, title: "write code", description: "finish the kanban project"},
		&Task{status: inProgress, title: "write code", description: "finish the kanban project"},
		&Task{status: inProgress, title: "write code", description: "finish the kanban project"},
	})

	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		&Task{status: done, title: "learn algebra", description: "repeat finite fields"},
		&Task{status: done, title: "learn algebra", description: "repeat finite fields"},
		&Task{status: done, title: "learn algebra", description: "repeat finite fields"},
	})

}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initLists(msg.Width, msg.Height)
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Previous()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveItemToNext
		case "n":
			models[modelMain] = m
			return models[modelForm].Update(msg)

		}
	case *Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "Loading ..."
	}
	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneView := m.lists[done].View()
	switch m.focused {
	case inProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			focusedStyle.Render(doneView),
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	}

}

/* FORM  MODEL */

type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	f := &Form{
		focused:     focused,
		title:       textinput.New(),
		description: textarea.New(),
	}
	f.title.Focus()
	return f
}

func (m *Form) CreateTask() tea.Msg {
	return NewTask(m.focused, m.title.Value(), m.description.Value())
}

func (m *Form) Init() tea.Cmd {
	return nil
}

func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			}
			models[modelForm] = m
			return models[modelMain], m.CreateTask
		}
	}

	var cmd tea.Cmd
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	}
	m.description, cmd = m.description.Update(msg)

	return m, cmd
}

func (m *Form) View() string {

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.title.View(),
		m.description.View(),
	)
}

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[modelMain]
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
