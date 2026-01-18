package main

// So far just initalized with the standard Bubbletea frame
// Added functionalities with my Noob Knowledge
// #HumbleKing

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel()) // Creates the actual running CLI
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	choices   []string         //items on the to-do lists
	cursor    int              //where the cursor pointing at
	selected  map[int]struct{} //which items are selected
	append    bool             //Append screen or not
	textinput textinput.Model  //Live User Input to add Items
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "New Item"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	var list []string
	list = ReadingPreviousList()
	return model{
		choices:   list,
		textinput: ti,
		selected:  make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.append {
			switch msg.String() {
			case "ctrl+c", "q":
				f, _ := os.Create("List.md")
				f.WriteString(FileFormating(m.choices))
				return m, tea.Quit
			case "a":
				m.append = true // Goes in to Append mode
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			case "d":
				m.choices = slices.Delete(m.choices, int(m.cursor), int(m.cursor)+1) // Deleting current selected Choice
			}
		} else {
			switch msg.String() {
			case "ctrl+l":
				m.append = false
			case "enter":
				m.choices = append(m.choices, m.textinput.Value())
				m.textinput.Reset()
				m.append = false
			}
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.append {
		return "Add a New Item: \n\n" + m.textinput.View() + "\n\nPress Ctrl+l to go back viewing the List"
		// "How the MD File looks like: \n" + FileFormating(m.choices)
	} else {
		s := "Your List:\n\n"

		for i, choice := range m.choices {

			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}

			s += fmt.Sprintf("%s [%s] %s \n", cursor, checked, choice)
		}

		s += "\n\nPress a to add a New Item\nPress d to delete the Selected Item\nPress q to quit.\nFiles Will be Saved in a extra MD File \n"

		return s
	}
}

// Formats the List in to the String which is used for the MD File
func FileFormating(list []string) string {
	// more efficient
	var s strings.Builder
	s.WriteString("List:")
	for _, item := range list {
		s.WriteString("\n" + item)
	}
	return s.String()
}

// It reads from already created MD File and returns items except for "List:"
func ReadingPreviousList() []string {
	body, err := os.ReadFile("List.md")
	if err != nil {
		fmt.Printf("A Failure wee oo wee oo %v", err)
	}
	var rendered []string
	rendered = strings.Split(string(body), "\n")
	return rendered[1:]
}
