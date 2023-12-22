package hosts

import (
	"KC-Checker/charm/threadPhase"
	"KC-Checker/helper"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strconv"
	"time"
)

var (
	style = lipgloss.NewStyle().Foreground(lipgloss.Color("#BE0101")).Width(threadPhase.GetWidth()).Align(lipgloss.Center)
	found = false
)

type waitingModel struct {
	spinner spinner.Model
}

func (m waitingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m waitingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	helper.GetProxiesFile("proxies.txt", true)

	if found {
		time.Sleep(time.Second * 3)
		return m, tea.Quit
	} else if helper.ProxySum > 0 {
		found = true
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m waitingModel) View() string {
	if found {
		return style.Copy().Foreground(lipgloss.Color("#01BE85")).Render("Found " +
			strconv.FormatFloat(helper.ProxySum, 'f', -1, 64) + " proxies")
	}

	return style.Render("It seems like you forgot to put proxies\nin proxies.txt\n"+
		"You don't have to close this tab\n") + style.Copy().Foreground(lipgloss.Color("#FFF")).Render(
		"\nIf you need proxies check out\nhttps://github.com/Kuucheen/KC-Scraper")
}

func WaitForProxies() {
	w := waitingModel{}
	if _, err := tea.NewProgram(w).Run(); err != nil {
		os.Exit(1)
	}
}
