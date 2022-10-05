package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	footerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	docStyle = lipgloss.NewStyle().Margin(1, 2)

	footer = "- ctrl+a: add book • ctrl+b: show books • ctrl+c/esc: exit "
)

func (b book) Title() string       { return b.title }
func (b book) Description() string { return b.desc }
func (b book) FilterValue() string { return b.title }

type book struct {
	title, desc string
}

type model struct {
	focusIndex int
	inputs     []textinput.Model
	window     int
	selected   map[int]struct{}
	list       list.Model
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) resetInputs() {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}

}

func (m *model) blurInputs() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func initialModel() model {

	existingLibrary := []list.Item{
		book{title: "Thinking, Fast and Slow", desc: "Daniel Kahneman"},
		book{title: "12 Rules for Life", desc: "Jordan B. Peterson"},
		book{title: "The Art of Seduction", desc: "Robert Greene"},
		book{title: "The 33 Strategies of War", desc: "Robert Greene"},
	}

	m := model{
		inputs:   make([]textinput.Model, 2),
		selected: make(map[int]struct{}),
		list:     list.New(existingLibrary, list.NewDefaultDelegate(), 0, 0),
	}

	m.list.Title = "Your Library"
	m.list.SetShowHelp(false)

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "Book"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Author"
			t.CharLimit = 64

		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, textinput.Blink)
	cmds = append(cmds, tea.EnterAltScreen)

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+a":
			m.inputs[0].Focus()
			m.window = 0

		case "ctrl+b":
			m.window = 1

		case "ctrl+d":
			if m.window == 1 {
				m.list.RemoveItem(m.list.Index())
			}

		case "enter", "up", "down":
			key := msg.String()

			cmds := make([]tea.Cmd, len(m.inputs))

			if m.window == 0 {

				if key == "up" {
					m.focusIndex--
				} else if key == "down" || key == "enter" {
					m.focusIndex++
				}

				if m.focusIndex < 0 {
					m.focusIndex = 2
				} else if m.focusIndex > 2 && key != "enter" {
					m.focusIndex = 0
				}

				if key == "enter" && m.focusIndex == 3 {

					m.window = 1
					m.focusIndex = 0

					if m.inputs[0].Value() != "" && m.inputs[1].Value() != "" {

						cmds = append(cmds, m.list.InsertItem(0, book{title: m.inputs[0].Value(), desc: m.inputs[1].Value()}))

					} else {
						return m, tea.Quit
					}

				}

				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}

			}

		}

	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	cmds = append(cmds, m.updateInputs(msg))

	return m, tea.Batch(cmds...)

}

func (m model) View() string {
	var b strings.Builder

	if m.window == 0 {

		b.WriteString("YourLibrary TUI\n\n")

		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}
		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n", *button)

	} else if m.window == 1 {

		m.resetInputs()
		m.blurInputs()

		return docStyle.Render(m.list.View())

	}
	b.WriteString(footerStyle(footer))

	return b.String()
}

func main() {

	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
