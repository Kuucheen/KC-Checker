package threadPhase

import (
	"KC-Checker/common"
	"KC-Checker/helper"
	"github.com/charmbracelet/bubbles/viewport"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	eliteCount       = 0
	anonymousCount   = 0
	transparentCount = 0
	banCheckCount    = 0

	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Width(GetWidth() / 2).Align(lipgloss.Center).Render
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Render

	centerStyle = lipgloss.NewStyle().Width(GetWidth() / 2).Align(lipgloss.Center).Render
)

type tickMsg time.Time

type model struct {
	//Displays how many % of total good proxies are elite,anonymous or transparent
	elite       progress.Model
	anonymous   progress.Model
	transparent progress.Model
	//Displays % of good proxies that passed the bancheck
	bancheck progress.Model
	//Displays % of progress made while checking
	percentage progress.Model

	viewport viewport.Model
	list     list.Model
}

var (
	threadPhase = true

	outputPath = ""
)

func RunBars() {
	items := []list.Item{
		item{title: "ip:port"},
		item{title: "type://ip:port"},
		item{title: "ip:port;ms"},
	}

	m := model{
		elite:       progress.New(progress.WithDefaultGradient()),
		anonymous:   progress.New(progress.WithDefaultGradient()),
		transparent: progress.New(progress.WithDefaultGradient()),
		bancheck:    progress.New(progress.WithDefaultGradient()),
		percentage:  progress.New(progress.WithDefaultGradient()),
		list:        list.New(items, list.NewDefaultDelegate(), 75, 0),
	}
	m.list.Title = "How do you want to save the proxies?"
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	m.list.SetShowHelp(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if helper.HasFinished {
		threadPhase = false
		SetStopTime()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if threadPhase {
			if msg.String() == "q" {
				threadPhase = false
				helper.StopThreads()
			}
		} else {
			switch msg.String() {
			case tea.KeyEnter.String():
				if common.DoBanCheck() {
					helper.Write(helper.ProxyMapFiltered, m.list.Index(), true)
				}
				outputPath = "\n\n" + helper.Write(helper.ProxyMap, m.list.Index(), false)
			case tea.KeyRight.String():
				m.list.CursorDown()
			case tea.KeyLeft.String():
				m.list.CursorUp()
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		width := GetWidth()/2 - 14
		m.elite.Width = width
		m.anonymous.Width = width
		m.transparent.Width = width
		m.bancheck.Width = width
		m.percentage.Width = width - 10
		m.list.SetWidth(msg.Width)
		return m, nil

	case tickMsg:
		sum := float64(getSum())

		eliteCount = helper.ProxyCountMap[3]
		anonymousCount = helper.ProxyCountMap[2]
		transparentCount = helper.ProxyCountMap[1]
		banCheckCount = helper.ProxyCountMap[-1]

		eliteCmd := m.elite.SetPercent(float64(eliteCount) / sum)
		anonymousCmd := m.anonymous.SetPercent(float64(anonymousCount) / sum)
		transparentCmd := m.transparent.SetPercent(float64(transparentCount) / sum)
		banCheckCmd := m.bancheck.SetPercent(float64(banCheckCount) / sum)
		percentageCmd := m.percentage.SetPercent(float64(getSumWithInvalid()) / helper.ProxySum)

		return m, tea.Batch(tickCmd(), eliteCmd, anonymousCmd, transparentCmd, banCheckCmd, percentageCmd)

	case progress.FrameMsg:
		eliteModel, eliteCmd := m.elite.Update(msg)
		m.elite = eliteModel.(progress.Model)

		anonymousModel, anonymousCmd := m.anonymous.Update(msg)
		m.anonymous = anonymousModel.(progress.Model)

		transparentModel, transparentCmd := m.transparent.Update(msg)
		m.transparent = transparentModel.(progress.Model)

		bancheckModel, bancheckCmd := m.bancheck.Update(msg)
		m.bancheck = bancheckModel.(progress.Model)

		percentageModel, percentageCmd := m.percentage.Update(msg)
		m.percentage = percentageModel.(progress.Model)

		return m, tea.Batch(eliteCmd, anonymousCmd, transparentCmd, bancheckCmd, percentageCmd)
	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := "\n\n"

	bars := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderRight(true).
		SetString(lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render("Elite"), "───────"+m.elite.View()) +
			pad + lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render("Anonymous"), "───"+m.anonymous.View()) +
			pad + lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render("Transparent"), "─"+m.transparent.View()) +
			pad + strings.Repeat("━", GetWidth()/2+2) +
			pad + lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render("Banchecked"), "──"+m.bancheck.View()))

	percentageBar := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		Width((GetWidth() - 10) / 2).
		BorderBottom(true).
		SetString("Progress  " + m.percentage.View())

	extraString := ""

	if threadPhase {
		extraString = helpStyle("Press q to stop")
	} else {
		extraString = centerStyle(m.list.View()) + "\n" + helpStyle("→ right • ← left • enter select")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, bars.String(),
		lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top, getStyledQueue(),
					getStyledInfo(eliteCount, anonymousCount, transparentCount)),
				percentageBar.String()), extraString), greenStyle(outputPath)),
	)
}

func (m model) renderLine(str string) string {
	title := titleStyle.Render(str)
	line := strings.Repeat("─", max(0, GetWidth()/2-lipgloss.Width(title)+2))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func getSum() int {
	sum := eliteCount + anonymousCount + transparentCount

	if sum == 0 {
		sum = 1
	}

	return sum
}

func getSumWithInvalid() int {
	return getSum() + helper.GetInvalid()
}
