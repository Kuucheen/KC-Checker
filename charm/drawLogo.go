package charm

import (
	"KC-Checker/charm/threadPhase"
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

var times int

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
		times += 2
		return m, tick
	}
	return m, nil
}

// View returns a string based on data in the logoModel. That string which will be
// rendered to the terminal.
func (m logoModel) View() string {
	logo := "  _  ______    ____ _               _             \n | |/ / ___|  / ___| |__   ___  ___| | _____ _ __ \n | ' / |     | |   | '_ \\ / _ \\/ __| |/ / _ \\ '__|\n | . \\ |___  | |___| | | |  __/ (__|   <  __/ |   \n |_|\\_\\____|  \\____|_| |_|\\___|\\___|_|\\_\\___|_|   \n \n\nby github.com/Kuucheen"

	width := threadPhase.GetWidth()

	var style = lipgloss.NewStyle().
		Width(width + 10).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#8271BB")).
		SetString(logo)

	var linestyle = style.Copy().
		Width(width).
		Align(lipgloss.Center).
		SetString(strings.Repeat("─", threadPhase.GetWidth()))

	str := style.Render() + "\n" + linestyle.Render() + "\n\n"

	return str
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Duration(100000 - times*900))
	return tickMsg{}
}
