package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

func main() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type (
	tickMsg struct{}
	errMsg  error
)

type model struct {
	textInput textinput.Model
	err       error
	window    int
	books     struct{ name []string }
	msg       string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Brave New World"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.window == 0 {
				if m.textInput.Value() == "" {
					m.msg = "⚠ Field can't be empty!\n"

				} else {
					m.books.name = append(m.books.name, m.textInput.Value())
					m.msg = fmt.Sprintf("[%s %s] added to your library.\n", randomEmoji(), m.books.name[len(m.books.name)-1])

					//m.textInput.Blur()
					m.textInput.Reset()
					return m, cmd
				}

			}

		}
		switch msg.String() {

		case "ctrl+a", "ctrl+b", "ctrl+x":
			m.msg = ""
			s := msg.String()

			// Add books
			if s == "ctrl+a" {
				m.window = 0
				m.textInput.Focus()
			}

			// Display books
			if s == "ctrl+b" {
				m.window = 2
				m.textInput.Blur()
			}

			// Main menu
			if s == "ctrl+x" {
				m.window = 3
				m.textInput.Blur()
			}

		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func randomEmoji() string {
	emojis := []rune("📔📕📖📗📘📙📓")
	return string(emojis[rand.Intn(len(emojis))])
}

func (m model) View() string {

	s := ""

	if m.window == 0 {
		s = fmt.Sprintf(
			"Enter the name of the book you want to add.\n\n%s\n\n",
			m.textInput.View(),
		)

	}

	if m.window == 2 {

		s = "Your library:\n\n"

		for _, book := range m.books.name {
			s += fmt.Sprintf("%s %s\n", randomEmoji(), book)
		}
	}

	if m.window == 3 {
		mb := m.books.name
		s = fmt.Sprintf("\n\n%d , %d\n\n", len(mb), cap(mb))
	}

	footer := helpStyle("\n- ctrl+a: add book • ctrl+b: show books • ctrl+c: exit\n")
	return s + m.msg + footer

}
