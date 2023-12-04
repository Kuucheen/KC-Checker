package charm

import (
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func DrawLogo() {
	// Initialize our program
	p := tea.NewProgram(logoModel(0))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// A logoModel can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.
type logoModel int

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m logoModel) Init() tea.Cmd {
	return tea.Batch(tick, tea.ClearScreen)
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated logoModel accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m logoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case tickMsg:
		if m >= 120 {
			return m, tea.Quit
		}
		m += 2
		return m, tick
	}
	return m, nil
}

// View returns a string based on data in the logoModel. That string which will be
// rendered to the terminal.
func (m logoModel) View() string {
	logo := "_  _ ____    ____ _  _ ____ ____ _  _ ____ ____ \n|_/  |    __ |    |__| |___ |    |_/  |___ |__/ \n| \\_ |___    |___ |  | |___ |___ | \\_ |___ |  \\ \n"

	var style = lipgloss.NewStyle().
		PaddingLeft(38).
		Foreground(lipgloss.Color("#758ECD")).
		SetString(logo)

	var linestyle = style.Copy().
		PaddingLeft(60 - int(m)/2).
		SetString(strings.Repeat("─", int(m)))

	str := style.Render() + "\n\n" + linestyle.Render() + "\n\n"

	return str
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 10)
	return tickMsg{}
}
