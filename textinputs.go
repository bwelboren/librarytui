package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"
	"strings"

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

	editMenu string
	footer   = "- ctrl+a: add book • ctrl+b: show books • ctrl+c/esc: exit "
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	window     int
	books      []string
	selected   map[int]struct{}
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
	m := model{
		inputs:   make([]textinput.Model, 2),
		books:    []string{},
		selected: make(map[int]struct{}),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Book"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 64

		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+a":
			// Put focus back to first index 'Book' when pressing ctrl+a
			m.inputs[0].Focus()
			m.window = 0

		case "ctrl+b":
			m.window = 1

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			cmds := make([]tea.Cmd, len(m.inputs))

			// Doesnt matter which window, we use the same keys

			var len_inputs int

			if m.window == 0 {
				len_inputs = len(m.inputs)
			} else if m.window == 1 {
				len_inputs = len(m.books)
			}

			if m.window == 0 {

				if s == "up" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}
				// When submitting
				if s == "enter" && m.focusIndex == len_inputs+1 {

					for _, text := range m.inputs {
						if text.Value() == "" {
							// Quit if given an empty value in either book or description
							return m, tea.Quit
						} else {
							m.books = append(m.books, text.Value())

						}
					}

					m.focusIndex = 0
					m.window = 1

				}

				// Handle textInput styling

				for i := 0; i <= len_inputs-1; i++ {
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

			if m.window == 1 {

				if s == "up" {
					if m.focusIndex > 0 {
						m.focusIndex--
					}
				}
				if s == "down" {
					if m.focusIndex < len_inputs-2 {
						m.focusIndex++
					}
				}

				if s == "enter" {
					_, ok := m.selected[m.focusIndex]
					if ok {
						delete(m.selected, m.focusIndex)
					} else {
						m.selected[m.focusIndex] = struct{}{}
					}
				}

			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	//books := make(map[string]string)

	if m.window == 0 {

		b.WriteString("Add a book.\n\n")

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
		b.WriteString("Your library.\n\n")

		editMenu = "• ctrl+d: delete •"

		x := 0
		for i := 0; i < len(m.books); i = i + 2 {
			cursor := " "
			if m.focusIndex == x {
				cursor = ">"
			}

			checked := " "
			if _, ok := m.selected[x]; ok {
				checked = "x"
			}
			x++

			s := fmt.Sprintf("%s [%s] Boek:%s\t\tBeschrijving:%s\n", cursor, checked, m.books[i], m.books[i+1])

			b.WriteString(s)

		}

	}

	b.WriteString(footerStyle(footer + editMenu))

	return b.String()
}

func main() {
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
