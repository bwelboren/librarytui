package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type model struct {
	err       error
	window    int
	Books     []string
	selected  map[int]struct{} // which to-do items are selected
	cursor    int
	spinner   spinner.Model
	textInput textinput.Model
}

type errMsg struct{ err error }

// e = errMsg
// e has err = Error
// error.Error
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func InitModel() model {
	sm := spinner.New()
	sm.Style = spinnerStyle
	sm.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Brave New World"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		spinner:   sm,
		err:       nil,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	//var cmds []tea.Cmd

	// Controleer de type van msg
	switch msg := msg.(type) {

	case errMsg:
		m.err = msg
		return m, tea.Quit

	// Als de type is dat er op een knop gedrukt is
	// Kijk welke knop zodra CTRL+C sluit dan af
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+b", "b":
			m.window = 1
		// Add book
		case "d":
			m.window = 2
		}

		// case "up", "k":
		// 	if m.cursor > 0 {
		// 		m.cursor--
		// 	}

		// // The "down" and "j" keys move the cursor down
		// case "down", "j":
		// 	if m.cursor < len(m.Books)-1 {
		// 		m.cursor++
		// 	}

		// // The "enter" key and the spacebar (a literal space) toggle
		// // the selected state for the item that the cursor is pointing at.
		// case "enter", " ":
		// 	_, ok := m.selected[m.cursor]
		// 	if ok {
		// 		delete(m.selected, m.cursor)
		// 	} else {
		// 		m.selected[m.cursor] = struct{}{}
		// 	}
		// }

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	// Als we geen ander bericht krijgen doe dan niks
	return m, cmd
}

func (m model) View() (s string) {

	s = "\n• Welcome •\n"

	s += fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), " ", textStyle("Running..."))
	footer := helpStyle("- d: add book • b: show books • q: exit\n")

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Laat gebruiker zien dat we bezig zijn

	if m.window == 1 {
		s += "\nYour books:\n"

		for i, choice := range m.Books {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := m.selected[i]; ok {
				checked = "x" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	if m.window == 2 {
		s += fmt.Sprintf(
			"What’s your favorite Book?\n\n%s\n\n%s",
			m.textInput.View(), "\n")
	}

	return "\n" + s + "\n\n" + footer

}

func main() {

	if err := tea.NewProgram(InitModel()).Start(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
