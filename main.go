package main

// So far just initalized with the standard Bubbletea frame
// Added functionalities with my Noob Knowledge
// #HumbleKing

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	input := Appending()
	p := tea.NewProgram(initialModel(input))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	choices  []string         //items on the to-do lists
	cursor   int              //where the cursor pointing at
	selected map[int]struct{} //which items are selected
	append   bool             //Append screen or not
	textarea textarea.Model
}

func initialModel(input []string) model {
	ta := textarea.New()
	ta.Placeholder = "Input: "
	list := slices.Concat(input, ReadingPreviousList())
	return model{
		choices:  list,
		textarea: ta,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			f, _ := os.Create("Hello.md")
			f.WriteString(FileFormating(m.choices))
			return m, tea.Quit
		case "a":
			m.append = true
		case "l":
			m.append = false
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
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.append {
		return "How the MD File looks like: \n" + FileFormating(m.choices)
	} else {
		s := "What you wanna enjoy\n\n"

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

		s += "\nPress q to quit.\n Files Will be Saved in a extra MD Folder \n"

		return s
	}
}

func Appending() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var inputs []string
	fmt.Print("Enter Your list \n")
	for 0 < 1 {
		scanner.Scan()
		if scanner.Text() == "q" || scanner.Text() == "" && len(inputs) != 0 {
			break
		} else if scanner.Text() == "q" || scanner.Text() == "" && len(inputs) == 0 {
			fmt.Printf("At least one non blank Input required\n")
		} else {
			inputs = append(inputs, scanner.Text())
		}
	}
	return inputs
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
	body, err := os.ReadFile("Hello.md")
	if err != nil {
		fmt.Printf("A Failure wee oo wee oo %v", err)
	}
	var rendered []string
	rendered = strings.Split(string(body), "\n")
	return rendered[1:]
}
