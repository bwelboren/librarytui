package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/charmbracelet/bubbles/spinner"

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
	spinner   spinner.Model
}

func initialModel() model {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	ti := textinput.New()
	ti.Placeholder = "..Book"
	ti.Focus()
	ti.CharLimit = 54
	ti.Width = 20

	return model{
		spinner:   sp,
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.window == 0 {
				if m.textInput.Value() == "" {
					m.msg = "⚠ Field can't be empty!\n"

				} else {
					m.books.name = append(m.books.name, m.textInput.Value())
					m.msg = fmt.Sprintf("%s %s added to your library.\n", randomEmoji(), m.books.name[len(m.books.name)-1])

					//m.textInput.Blur()
					m.textInput.Reset()
					return m, cmd
				}
			}

		case "ctrl+a", "ctrl+b", "ctrl+x":
			s := msg.String()

			m.msg = ""
			m.textInput.Blur()

			// Add books
			if s == "ctrl+a" {
				m.window = 0
				m.textInput.Focus()
			}

			// Display books
			if s == "ctrl+b" {
				m.window = 1
			}


		}

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

	s := m.spinner.View() + " "

	if m.window == 0 {
		s += fmt.Sprintf(
			"Enter the name of the book you want to add.\n\n%s\n\n",
			m.textInput.View(),
		)

	}

	if m.window == 1 {
		s += "Your library:\n\n"
		for _, book := range m.books.name {
			s += fmt.Sprintf("📘 %s\n", book)
		}
	}

	footer := helpStyle("\n- ctrl+a: add book • ctrl+b: show books • ctrl+c/esc: exit\n")
	return s + m.msg + footer

}
