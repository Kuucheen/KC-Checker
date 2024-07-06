package charm

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var index = 0
var finished = false

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("KC-Checker - github.com/Kuucheen")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			finished = true
			index = m.list.Index()
			return m, tea.Quit
		} else if msg.String() == "ctrl+c" {
			os.Exit(1)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-10)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if finished {
		return ""
	}
	return docStyle.Render(m.list.View())
}

func GetProxyType() int {

	items := []list.Item{
		item{"HTTP", "for various web applications"},
		item{"HTTPS", "for various web applications"},
		item{title: "SOCKS4", desc: "for various app applications"},
		item{title: "SOCKS5", desc: "for various app applications"},
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "What type of proxies do you want to check?"
	//m.list.SetStatusBarItemName("Type", "Types")
	m.list.SetShowStatusBar(false)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return index
}
