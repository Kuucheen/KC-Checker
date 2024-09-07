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

type logoModel int

var times int

func (m logoModel) Init() tea.Cmd {
	return tea.Batch(tick, tea.ClearScreen, tea.SetWindowTitle("KC-Checker - github.com/Kuucheen"))
}

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

func (m logoModel) View() string {
	logo := " _     _ _______    _______ _                 _             \n(_)   | (_______)  (_______) |               | |            \n    _____| |_          _      | |__  _____  ____| |  _ _____  ____ \n   |  _   _) |        | |     |  _ \\| ___ |/ ___) |_/ ) ___ |/ ___)\n| |  \\ \\| |_____   | |_____| | | | ____( (___|  _ (| ____| |\n|_|   \\_)\\______)   \\______)_| |_|_____)\\____)_| \\_)_____)_|\n\n"
	logo += "by github.com/Kuucheen"

	width := threadPhase.GetWidth()

	var style = lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#6C44BB")).
		SetString(logo)

	var linestyle = style.
		Width(width).
		Align(lipgloss.Center).
		SetString(strings.Repeat("─", threadPhase.GetWidth()))

	str := style.Render() + "\n" + linestyle.Render() + "\n"

	return str
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Duration(100000 - times*900))
	return tickMsg{}
}
